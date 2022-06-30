package memotest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type PlayerId struct { val int }

func (id PlayerId) str() string { return strconv.Itoa(id.val) }
func (id* PlayerId) inc() { (*id).val++ }
func PlayerIdFromStr(s string) (PlayerId, error) {
	val, err := strconv.Atoi(s)
	return PlayerId{val}, err
}

type Player 	struct {
	Loop		*Loop
	Id    		PlayerId
	Name		string
	Extra 		any
	Playing		*Game
}

func NewPlayer(id PlayerId, name string, extra any) *Player {
	if(0==len(name)) {
		name = "Anon"+id.str()
	}
	player 		:= 	Player{nil, id, name, extra, nil}
	player.Loop =	NewLoop(&player)
	return 			&player
}

func PlayerAsync[T any](player *Player,resp RetWithError[T], fn LoopFn) RetWithError[T] {
	if (player == nil) {
		go func() {
			ret := WithError[T]{}
			ret.err = errors.New("Null player")
			resp.SendAndClose( ret )
		} ()
	} else {
		player.Loop.Async(fn)
	}
	return resp
}


func (player *Player) GetId() RetWithError[PlayerId] {
	resp := NewRetWithError[PlayerId]()
	// Podría ser player.Loop.async pero realmente no hace falta bloquear
	// Pero sí debe ser un goroutine porque quien la llama es quien espera
	// el canal, y se daría un bloqueo.
	go func() {
		if(player == nil) {
			resp.SendNewAndClose( PlayerId{0}, errors.New("Null player") )
		} else {
			resp.SendNewAndClose( player.Id, nil )
		}
	} ()
	return resp
}

func (player *Player) IsValid() chan bool {
	return AsyncVal(player != nil)
}

/** @brief Crea un juego si no está en ninguno, y se une. **/
/** \todo Debería tener referencia a Players, que tiene a Games **/
func (player *Player) NewGame(games *Games, config GameConfig, extra any) RetWithError[*Game] {
	resp := NewRetWithError[*Game]()
	return PlayerAsync( player, resp, func(loop *Loop) {
		// Si está en un juego, no puede crear uno.
		if (player.Playing != nil) {
			previous := <- player.Playing.GetId()
			txt := fmt.Sprintf(`Already playing %v`, previous.val.str() )
			resp.SendNewAndClose( nil, errors.New(txt) )
			return
		}
		intent := <- games.NewGame(config, extra)
		if(intent.err == nil) { 		// Si todo bien
			game := intent.val
			intent = <- game.Join(player)	// Intenta unirse
			if(intent.err == nil) {
				player.Playing = game
			} else {
				game.Kill("Can't join creator")
			}
		}
		resp.SendAndClose( intent )
	})
}

func (player *Player) JoinGame(game *Game) RetWithError[*Game] {
	resp := NewRetWithError[*Game]()
	return PlayerAsync( player, resp, func(loop *Loop) {
		// Si está en un juego, no puede unirse a uno
		if (player.Playing != nil) {
			previous := <- player.Playing.GetId()
			txt := fmt.Sprintf(`Already playing %v`, previous.val.str() )
			resp.SendNewAndClose( nil, errors.New(txt) )
			return
		}
		intent := <- game.Join(player)
		resp <- intent
		player.Playing = intent.val // el juego, o nil en caso de error
	})
}

func (player *Player) Show() chan string {
	resp := NewRetWithError[SStr]()
	return ErrToJson( PlayerAsync[SStr]( player, resp, func(loop *Loop) {
		stream := make(chan string)
		/** \todo gameId (null si game == nil) **/
		plId := player.Id
		name, err := json.Marshal(player.Name)
		go func() { // Fuera de bucle player
			stream <- `{"gameId":null`
			stream <- `,"playerId":` + plId.str()
			stream <- `,"name":`
			
			if( err == nil) {
				stream <- string(name)
			} else {
				stream <- "null"
			}
			stream <- "}"
			close(stream)	
		} () // Fin: Fuera de bucle player
		resp.SendNewAndClose(stream,nil)
	}) )
}