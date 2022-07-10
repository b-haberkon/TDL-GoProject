package memotest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
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
	Selected	*Piece
	Owned		[]*Piece
}

func NewPlayer(id PlayerId, name string, extra any) *Player {
	if(0==len(name)) {
		name = "Anon"+id.str()
	}
	player 		:= 	Player{nil, id, name, extra, nil, nil, make([]*Piece, 2)}
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
		// Si ya está en el mismo juego, devuelve el mismo juego
		if (player.Playing == game) {
			resp.SendNewAndClose(game, nil)
			return
		}
		// Si está en otro juego, no puede unirse a uno
		if (player.Playing != nil) {
			previous := <- player.Playing.GetId()
			txt := fmt.Sprintf(`Already playing %v`, previous.val.str() )
			resp.SendNewAndClose( nil, errors.New(txt) )
			return
		}
		intent := <- game.Join(player)
		resp <- intent
		if(intent.err == nil) {
			player.Playing = intent.val // el juego, o nil en caso de error
		}		
	})
}

func (player *Player) ShowWith(extra string) chan string {
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
			if(extra != "") {
				stream <- `,` + extra
			}
			stream <- "}"
			close(stream)	
		} () // Fin: Fuera de bucle player
		resp.SendNewAndClose(stream,nil)
	}) )
}

func (player *Player) selectPiece(gameId GameId,pieceId PieceId) RetWithError[*MoveResult] {
	resp := NewRetWithError[*MoveResult]()
	return PlayerAsync( player, resp, func(loop *Loop) {
		// Ver que el gameId coincida
		myGame := <- player.Playing.GetId()
		if (myGame.err != nil) || (myGame.val != gameId) {
			msg := fmt.Sprintf("You are playing %v, not %v", myGame.val.str(), gameId.str())
			resp.SendNewAndClose(nil,errors.New(msg))
			return
		}
		// Obtener pieza
		pieceIntent := <- player.Playing.getPiece(pieceId)
		if(pieceIntent.err != nil) {
			resp.SendNewAndClose(nil,pieceIntent.err)
			return
		}
		/** \todo La pieza podría devolver un resultado, según si estaba ocupada,
		 * no existe, etc., y devolver eso al usuario. **/
		/** \todo El jugador podría tener un estado para no elegir muchas piezas rápidamente;
		 * además, mientras espera una respuesta podría devolver el error de "seleccionando".
		 **/
		piece := pieceIntent.val;
		//var intent WithError[*MoveResult]
		intent := <- piece.SelectOrPair(player, player.Id, player.Selected)
		if(intent.err == nil) {
			switch( intent.val.Status ) {
			case Match:
				player.Owned = append(player.Owned, piece, player.Selected)
				player.Selected = nil
			case Unmatch:
				player.Selected = nil
			default:
				player.Selected = intent.val.Pieces[0]
			} // switch
		} // if (! intent.err )
		fmt.Printf("** DEBUG5: Selected=%v\n",player.Selected)
		resp <- intent
	}) // async
}

func (player *Player) End(wg *sync.WaitGroup) {
	if(player == nil) {
		go func() { wg.Done() } ()
	} else {
		player.Loop.Async(func(loop *Loop) {
			if(player.Playing != nil) {
				player.Playing = nil
			}
			wg.Done();
		})
	}
}
