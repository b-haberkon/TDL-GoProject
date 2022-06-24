package memotest

import(
	"errors"
	//"fmt"
)

type Games	struct {
	Loop 	*Loop
	Ids   	GameId
	Extra 	any
	All		map[GameId]*Game
}

func NewGames(extra any) *Games {
	games 		:=	Games{nil, GameId{0}, extra, make(map[GameId]*Game)}
	games.Loop	=	NewLoop(&games)
	return			&games
}

func GamesAsync[T any](games *Games,resp RetWithError[T], fn LoopFn) RetWithError[T] {
	if (games == nil) {
		go func() {
			ret := WithError[T]{}
			ret.err = errors.New("Null games")
			resp.SendAndClose( ret )
		} ()
	} else {
		games.Loop.Async(fn)
	}
	return resp
}

func (games *Games) NewGame(config GameConfig, extra any) RetWithError[*Game] {
	resp := make(RetWithError[*Game])
	return GamesAsync( games, resp, func(loop *Loop) {
		games.Ids.inc()
		id := games.Ids
		go func() { // Ya no necesita bloquear games
			resp.SendNewAndClose( NewGame(id, config, extra), nil)
		} ()
	})
}

func (games *Games) GetById (id GameId) RetWithError[*Game] {
	resp := NewRetWithError[*Game]()
	return GamesAsync( games, resp, func(loop *Loop) {
		game := games.All[id]
		if( <- game.IsValid() ) {
			resp.SendNewAndClose(game,nil)
		} else {
			resp.SendNewAndClose(nil,errors.New("Invalid game"))
		}
	})
}