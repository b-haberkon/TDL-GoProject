package memotest

import(
	"errors"
	//"fmt"
)

type Players struct {
	Loop	*Loop
	Ids   	PlayerId
	Extra 	any
	Games 	*Games
	All		map[PlayerId]*Player
}

func NewPlayers(games *Games, extra any) *Players {
	players 		:=	Players{nil, PlayerId{0}, extra, games, make(map[PlayerId]*Player)}
	players.Loop	=	NewLoop(&players)
	return 				&players
}

func(players *Players) NewPlayer(name string, extra any) RetWithError[*Player] {
	resp := make(RetWithError[*Player])
	PlayersAsync( players, resp, func(loop *Loop) {
		players.Ids.inc()
		id := players.Ids
		go func() { // Por ahora desbloquea players
			player := NewPlayer(id, name, extra)
			// Bloqueará hasta que players vuelva a tomar control,
			// agregándolo al map y enviando la respuesta.
			players.SetPlayer(id, player, resp)
		} ()
	})

	return resp
}

func PlayersAsync[T any](players *Players,resp RetWithError[T], fn LoopFn) RetWithError[T] {
	if (players == nil) {
		go func() {
			ret := WithError[T]{}
			ret.err = errors.New("Null players")
			resp.SendAndClose( ret )
		} ()
	} else {
		players.Loop.Async(fn)
	}
	return resp
}

func (players *Players) SetPlayer (id PlayerId, player *Player, resp RetWithError[*Player]) RetWithError[*Player] {
	resp.SendNewAndClose( player, nil )
	return PlayersAsync( players, resp, func(loop *Loop) {
		players.All[id] = player
	})
}

func (players *Players) GetById (id PlayerId) RetWithError[*Player] {
	resp := NewRetWithError[*Player]()
	return PlayersAsync( players, resp, func(loop *Loop) {
		player := players.All[id]
		if( <- player.IsValid() ) {
			resp.SendNewAndClose(player,nil)
		} else {
			resp.SendNewAndClose(nil,errors.New("Invalid player"))
		}
	})
}