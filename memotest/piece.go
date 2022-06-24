package memotest

import (
	//"fmt"
	"encoding/json"
	"strconv"
)
type PieceId int8
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
    piece       :=  &Piece{nil, id, row, col, *symbol, Hidden, 0}
    piece.Loop  =   NewLoop(piece)
    return          piece
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
				stream <- `,"Id":`  + strconv.Itoa(int( piece.Id  ))
				stream <- `,"Row":` + strconv.Itoa(int( piece.Row ))
				stream <- `,"Col":` + strconv.Itoa(int( piece.Col ))
				stream <- `}`
				close(stream)
				} ()
		})
	}
	return stream
}