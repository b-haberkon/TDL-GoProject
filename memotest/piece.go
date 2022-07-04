package memotest

import (
	//"fmt"
	"encoding/json"
	"errors"
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
    State   PieceState  `json:"State"`
	SelBy   PlayerId	`json:"SelBy"`
}

var emptySymbol Symbol = Symbol{"", -1}

func NewPiece(id PieceId, row uint8, col uint8, symbol *Symbol) *Piece {
	if(symbol == nil) {
		symbol = &emptySymbol;
	}
    piece       :=  &Piece{nil, id, row, col, *symbol, Hidden, PlayerId{0}}
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
				stream <- `}`
				close(stream)
				} ()
		})
	}
	return stream
}

func (piece *Piece) Select(player *Player, playerId PlayerId) RetWithError[*MoveResult] {
	resp := NewRetWithError[*MoveResult]()
	ret := WithError[*MoveResult]{NewMoveResult(Inexistent), nil}
	return PieceAsync( piece, resp, func(loop *Loop) {
		defer func() { resp.SendAndClose(ret)} () // Pase lo que pase, enviar una respuesta
		ret.val.Pieces = append(ret.val.Pieces, piece)
		switch piece.State {
		case Hidden:
			piece.toState( Selected, playerId )
			ret.val.Status = Selection;
		case Removed:
			ret.val.Status = Inexistent
		default:
			ret.val.Status = Blocked
		}
	})
}

func (piece *Piece) Pair(player *Player, playerId PlayerId, another *Piece) RetWithError[*MoveResult] {
	/** \todo Evitar interbloqueo si ambas se seleccionan a mensajean a la vez. **/
	resp := NewRetWithError[*MoveResult]()
	ret := WithError[*MoveResult]{NewMoveResult(Blocked), nil}
	return PieceAsync( piece, resp, func(loop *Loop) {
		defer resp.SendAndClose(ret)
		if (piece.State == Hidden) && (piece != another) { // La segunda pieza debe estar oculta
			ret.val.Status = another.isMatch(piece.Symbol, playerId) // ¿La primera pieza coincide?
			if (ret.val.Status == Match) {
				piece.toState(Matched, playerId);
			} else if(ret.val.Status == Unmatch) {
				piece.toState(Unmatched, playerId);
			}
		}
	})
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
/** Llamada sólo por otra pieza. **/
func (piece *Piece) isMatch(symbol Symbol, playerId PlayerId) MoveResultStatus {
	if(piece == nil) {
		return Inexistent
	}
	resp := make(chan MoveResultStatus)
	piece.Loop.Async( func(loop *Loop) {
		if (piece.SelBy != PlayerId{0}) {
			resp <- Blocked
			return
		} else if (piece.Symbol.Pair == symbol.Pair) {
			resp <- Match
			piece.toState(Matched, playerId)
		} else {
			resp <- Unmatch
			piece.toState(Unmatched, playerId)
		}
	})
	return <- resp
}

/** Ejecutada dentro del bucle. **/
func (piece *Piece) toState(state PieceState, playerId PlayerId) {
	if(piece == nil) {
		return
	}
	piece.State = state
	piece.SelBy = playerId
	piece.Loop.Async(func(loop *Loop) {
		to, _ := timeout(stateExpiration[state])
		if (<- to) && (piece.State == state) && ( piece.SelBy == playerId ) {
			piece.State = stateAfterExpiration[state]
			piece.SelBy = PlayerId{0}
			/** ¿Notificación? **/
		}
	})
}