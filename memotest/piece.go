package memotest

import (
	//"fmt"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)
type PieceId  struct { val int }
func (id PieceId) str() string { return strconv.Itoa(id.val) }
func (id* PieceId) inc() { (*id).val++ }
func PieceIdFromStr(s string) (PieceId, error) {
	val, err := strconv.Atoi(s)
	return PieceId{val}, err
}

type PieceState int8

const (
    Hidden PieceState = iota    // Estado inicial, la pieza no es visible.
    Selected                    // La pieza fue seleccionada por un jugador (visible)
    Matched                     // La pieza fue emparejada con otra por un jugador
    Unmatched                   // La pieza no pudo ser emparejada con otra por un jugador
    Removed                     // La pieza fue removida al ser emparejada
)
var pieceStateToText = map[PieceState]string {
    Hidden		:	"Hidden",
    Selected	:	"Selected",
    Matched		:	"Matched",
    Unmatched	:	"Unmatched",
    Removed		:	"Removed",
}
func (st PieceState) str() string { return pieceStateToText[st] }

var pieceStateVisibility = map[PieceState]bool {
    Hidden		:	false,
    Selected	:	true,
	Matched		:	true,
	Unmatched	:	true,
	Removed		:	false,
}
func (st PieceState) isVisible() bool { return pieceStateVisibility[st] }

type Piece struct {
    Loop    *Loop
    Id      PieceId     `json:"id"`
	Row		uint8		`json:"Row"`
	Col		uint8		`json:"Col"`
    Symbol  Symbol      
	Game	*Game		`json:"Game"`
    State   PieceState  `json:"State"`
	SelBy   PlayerId	`json:"SelBy"`
	cancel  chan bool
}

var emptySymbol Symbol = Symbol{"", 0}

func NewPiece(id PieceId, row uint8, col uint8, symbol *Symbol, game *Game) *Piece {
	if(symbol == nil) {
		symbol = &emptySymbol;
	}
    piece       :=  &Piece{nil, id, row, col, *symbol, game, Hidden, PlayerId{0}, nil}
    piece.Loop  =   NewLoop(piece)
    return          piece
}


func NewEmptyPiece(id PieceId, row uint8, col uint8) *Piece {
    piece       :=  &Piece{nil, id, row, col, Symbol{"",0}, nil, Removed, PlayerId{0}, nil}
    piece.Loop  =   NewLoop(piece)
    return          piece
}

func PieceAsync[T any](piece *Piece,resp RetWithError[T], fn LoopFn) RetWithError[T] {
	if (piece == nil) {
		go func() {
			ret := WithError[T]{}
			ret.err = errors.New("Null piece")
			resp.SendAndClose( ret )
		} ()
	} else {
		piece.Loop.Async(fn)
	}
	return resp
}


func (piece *Piece) Show() chan string {
	stream := make(chan string)
	if(piece == nil) {
		go func() {
			stream <- `null`
			close(stream)
		} ()
	} else {
		piece.Loop.Async(func(loop *Loop){
			state := piece.State
			text, err := json.Marshal(piece.Symbol.Text)

			go func() {	// El resto no cambia => no bloquea
				if (err != nil) || !state.isVisible() {
					text = []byte(`null`)
				}
				stream <- `{"State":"` + state.str() + `"`
				stream <- `,"Text":`  + string(text)
				stream <- `,"Id":`  + piece.Id.str()
				stream <- `,"Row":` + strconv.Itoa(int( piece.Row ))
				stream <- `,"Col":` + strconv.Itoa(int( piece.Col ))
				stream <- `,"SelBy":` + piece.SelBy.str()
				stream <- `}`
				close(stream)
				} ()
		})
	}
	return stream
}

var stateExpiration = map[PieceState]time.Duration {
    Hidden		:	0,
    Selected	:	10 * time.Second,
    Matched		:	4 * time.Second,
    Unmatched	:	2 * time.Second,
    Removed		:	0,
}
var stateAfterExpiration = map[PieceState]PieceState {
    Hidden		:	Hidden,
    Selected	:	Hidden,
    Matched		:	Removed,
    Unmatched	:	Hidden,
    Removed		:	Removed,
}

/** Llamada sólo por la segunda pieza. **/
func (piece *Piece) isMatch(symbol Symbol, player *Player, playerId PlayerId) MoveResultStatus {
	if(piece == nil) {
		return Selection
	}
	resp := make(chan MoveResultStatus)
	piece.Loop.Async( func(loop *Loop) {
		available := (piece.State == Selected) && (piece.SelBy == playerId)
		if (! available ) {
			resp <- Selection
			return
		}
		
		if (piece.Symbol.Pair == symbol.Pair) {
			resp <- Match
			piece.toState(Matched, player, playerId)
		} else {
			resp <- Unmatch
			piece.toState(Unmatched, player, playerId)
		}
	})
	return <- RespOrTimeout(resp, 1*time.Second, func() MoveResultStatus {
		fmt.Printf("piece %v isMatch player %v timeout",piece.Id.str(),playerId.str())
		return Selection // Unavailable
	}, nil)
}

/** Ejecutada dentro del bucle. **/
func (piece *Piece) toState(state PieceState, player *Player, playerId PlayerId) {
	if(piece == nil) {
		return
	}
	piece.State = state
	fmt.Printf("Piece %v to state %v…\n",piece.Id,state)
	piece.SelBy = playerId
	msExpiration := stateExpiration[state]
	if( msExpiration == 0) {
		piece.cancel = nil
	} else {
		// Programa un timeout
		expiration, _, cancel := timeout(msExpiration)
		// cancel quedará como closure, actualiza piece
		piece.cancel = cancel
		go func() {
			res := <- expiration
			piece.Loop.Async(func(*Loop) {
				// Expiró y no cambió el timer (closure coincide con piece)
				if (res) && (piece.cancel == cancel) {
					piece.State = stateAfterExpiration[state]
					piece.SelBy = PlayerId{0}
					fmt.Printf("Piece %v back to state %v…\n",piece.Id,piece.State)
					if(piece.State == Removed) {
						piece.Game.PieceRemoved(piece, player, 1)
					}
				}
			})
		} ()
	}
}

/**
 * Intenta emparejar la pieza actual con previous, si previous está seleccionada.
 * O si no, sólo intenta seleccionar la pieza actual.
 * Puede devolver:
 *   Status  | Pieces          | Cuándo
 * Selection | actual,previous | La pieza previa no era válida o no está seleccionada por el jugador. Pero la actual sí.
 * --otro--  | previous,actual | La pieza actual no era válida o no está hidden. La previa tal vez.
 * Si la pieza previa estaba seleccionada por el jugador, y la pieza actual hidden:
 * Match     | ambas           | Coinciden además en par.
 * Unmatch   | ambas           | Pero no coinciden en par.
 */
func (piece *Piece) SelectOrPair(player *Player, playerId PlayerId, previous *Piece) RetWithError[*MoveResult] {
	resp := NewRetWithError[*MoveResult]()
	ret := WithError[*MoveResult]{NewMoveResult(Blocked), nil}
	return PieceAsync( piece, resp, func(loop *Loop) {
		defer resp.SendAndClose(ret)
		// Si la actual no es válida, devuelve primero la otra
		// para mantenerla como selección previa (sea o no buena)
		fmt.Printf("*** DEBUG1: piece=%v,prev=%v,ret=%v\n",piece,previous,ret)
		if (piece.State != Hidden ) || (piece.SelBy == playerId) {
			if (piece.State == Removed) {
				ret.val.Status = Inexistent
			} else {
				ret.val.Status = Blocked
			}
     		ret.val.Pieces = append(ret.val.Pieces, previous, piece)
			 fmt.Printf("*** DEBUG2: piece=%v,prev=%v,ret=%v\n",piece,previous,ret)
			 return
		}
		// En cualquier otro caso, la actual va primero;
		// pero devuelve ambas para que se actualicen ambos estados.
		ret.val.Pieces = append(ret.val.Pieces, piece, previous)
		fmt.Printf("*** DEBUG3: piece=%v,prev=%v,ret=%v,%v\n",piece,previous,ret.val.Status,ret.val.Pieces)

		// Estando ésta disponible para matchear, intenta hacerlo con la previa
		ret.val.Status = previous.isMatch(piece.Symbol, player, playerId)
		fmt.Printf("*** DEBUG4a: match status=%v\n",ret.val.Status)
		switch( ret.val.Status ) {
		case Match:		// previous es seleccionable, y hubo coincidencia
			piece.toState(Matched, player, playerId)	// Acompaña a la otra en estado
		case Unmatch: 	// previous es seleccionable, y no hubo coincidencia
			piece.toState(Unmatched, player, playerId)	// Acompaña a la otra en estado
		default:		// En otro caso, previous no se modifica, pero piece es válida
			piece.toState( Selected, player, playerId )	// Entonces se selecciona piece
			ret.val.Status = Selection;			// Y siempre será una selección
		}
		fmt.Printf("*** DEBUG4b: piece=%v,prev=%v,ret=%v,%v\n",piece,previous,ret.val.Status,ret.val.Pieces)
	})
}

func (piece *Piece) Die() {
	if(piece == nil) {
		return
	}
	piece.Loop.Async(func(loop *Loop) {
		loop.Close()
	})
}