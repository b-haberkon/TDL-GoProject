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
}

func NewPlayers(games *Games, extra any) *Players {
	players 		:=	Players{nil, 0, extra, games}
	players.Loop	=	NewLoop(&players)
	return 				&players
}

func(players *Players) NewPlayer(name string, extra any) RetWithError[*Player] {
	resp := make(RetWithError[*Player])
	if( nil == players) {
		return resp.SendNewAndClose(nil, errors.New("Null player"))
	}
	//ret := WithError[*Player]{nil, nil}
	players.Loop.Async(func(loop *Loop) {
		players.Ids++
		id := players.Ids		
		go func() { // A partir de ahora ya no necesita bloquear Players
			//defer resp.SendAndClose(ret)
			// player := <- NewPlayer(id, name)
			//ret.val = player
			resp.SendNewAndClose( NewPlayer(id, name, extra), nil )
		} ()
	})

	return resp
}