package memotest

import (
	//"fmt"
	"encoding/json"
	"errors"
	"strconv"
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
    Paired                      // La pieza fue emparejada con otra por un jugador
    Removed                     // La pieza fue removida al ser emparejada
)
var pieceStateToText = map[PieceState]string {
    Hidden		:	"Hidden",
    Selected	:	"Selected",
    Paired		:	"Paired",
    Removed		:	"Removed",
}
func (st PieceState) str() string { return pieceStateToText[st] }

var pieceStateVisibility = map[PieceState]bool {
    Hidden		:	false,
    Selected	:	true,
    Paired		:	true,
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
	ret := WithError[*MoveResult]{&MoveResult{Inexistent,make([]*Piece,2)}, nil}
	return PieceAsync( piece, resp, func(loop *Loop) {
		defer func() { resp.SendAndClose(ret)} () // Pase lo que pase, enviar una respuesta
		ret.val.Pieces = append(ret.val.Pieces, piece)
		switch piece.State {
		case Hidden:
			piece.State = Selected
			piece.SelBy = playerId
			ret.val.Status = Selection;
		default:
			ret.val.Status = Blocked
		}
	})
}

func (piece *Piece) Pair(player *Player, playerId PlayerId, another *Piece) RetWithError[*MoveResult] {
	/** \todo Usar timeout para evitar interbloqueo **/
	resp := NewRetWithError[*MoveResult]()
	return PieceAsync( piece, resp, func(loop *Loop) {
		resp.SendNewAndClose(nil, errors.New("Unimplemented"))
	})
}