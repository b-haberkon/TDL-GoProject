package memotest

import (
	"errors"
	"fmt"
	"strconv"
)

type GameId struct { val int }

func (id GameId) str() string { return fmt.Sprintf("%v",id.val)}
func (id* GameId) inc() { (*id).val++ }

type GameStatus uint8

const (
	Waiting GameStatus = iota
	Playing
	Ended
	Dead
)
var gameStatusToText = map[GameStatus]string {
	Waiting	:	"Waiting",
	Playing	:	"Playing",
	Ended	:	"Ended",
	Dead	:	"Dead",
}
func (st GameStatus) str() string { return gameStatusToText[st] }

type GameConfig struct {
	Rows uint8
	Cols uint8
	Syms []*Symbol
	PMin uint8
	PMax uint8
}

type Game struct {
	Loop	*Loop
	Id		GameId
	Config	GameConfig
	Extra	any
	Players	[]*Player
	Pieces  map[PieceId]*Piece
	Status  GameStatus
}

func limitate[T int | uint | uint8 | ~string ](ptr *T, min T, max T) {
	if( (*ptr) < min) {
		(*ptr) = min
	} else
	if( (*ptr) > max) {
		(*ptr) = max
	}
}
func min[T int](a,b int) int {
	if(a<b) {
		return a
	}
	return b
}

func max[T int](a,b int) int {
	if(a>b) {
		return a
	}
	return b
}

func shuffleSymbols(org []*Symbol, amount int) []*Symbol {
	/** \todo Agrupar de a pares y otras cosas. **/
	dest := make([]*Symbol, amount)
	l := len(org)
	copy(dest, org[0:min(amount,l)])
	return dest
}

func NewGame(id GameId, config GameConfig, extra any) *Game {
	game 		:= 	Game{
		Loop	:	nil,
		Id		:	id,
		Config	:	config,
		Extra	:	extra,
		Players	:	make([]*Player, 0, 2),
		Pieces	:	make(map[PieceId]*Piece),
		Status  :   Waiting }
	game.Loop	=	NewLoop(&game)
	// Usa WaitTurn para no salir de NewGame hasta que
	// el bucle haya leído la función (para así garantizar
	// que sea la primera)
	game.Loop.WaitTurn(func(loop *Loop) {
		limitate[uint8] (	&game.Config.Rows,	2, 16)
		limitate[uint8] (	&game.Config.Cols,	2, 16)
		limitate[uint8] (	&game.Config.PMin,	1, 16)
		limitate[uint8] (	&game.Config.PMax,	game.Config.PMin, 16)
		nPieces := int(game.Config.Rows) * int(game.Config.Cols)
		if( (nPieces % 2) > 0) {
			nPieces--
		}
		dest := shuffleSymbols(game.Config.Syms, nPieces)
		var row, col uint8
		var i PieceId = 0
		for row=0; row<game.Config.Rows; row++ {
			for col=0; col<game.Config.Cols; col++ {
				pos := int(i) % len(dest)
				i++
				game.Pieces[i] = NewPiece(i, row, col, dest[pos])
			}
		}
	})
	return &game
}

func (game *Game) GetId() RetWithError[GameId] {
	resp := NewRetWithError[GameId]()
	go func() {
		if(nil == game) {
			resp.SendNewAndClose( GameId{0}, errors.New("Null game") )
		} else {
			resp.SendNewAndClose( game.Id, nil )
		}
	} ()
	return resp
}

func (game *Game) IsValid() chan bool {
	resp := make(chan bool)
	go func() {
		resp <- (game == nil)
		close(resp)
	} ()
	return resp
}

func (game *Game) Join(player *Player) RetWithError[*Game] {
	resp := NewRetWithError[*Game]()
	if(game == nil) {
		return resp.SendNewAndClose(nil, errors.New("Invalid game") )
	}
	if(player == nil) {
		return resp.SendNewAndClose(nil, errors.New("Invalid player") )
	}
	game.Loop.Async(func(loop *Loop){
		/** \todo game.IsValid() player.IsValdi() **/
		/** \todo game not end **/
		var amount uint8 = uint8( len(game.Players) )
		if(amount >= game.Config.PMax) {
			txt := fmt.Sprintf("Game is full, there are %v, max: %v",
				amount, game.Config.PMax)
			resp.SendNewAndClose(nil, errors.New(txt))
			return
		}

		if(game.Config.PMin == amount + 1) && (Waiting == game.Status) {
			/** \todo Notificar, etc.*/
			game.Status = Playing
		}

		game.Players = append(game.Players, player)
		resp.SendNewAndClose(game, nil)
		fmt.Printf("Game %v join player %v\n",game.Id,player.Id)
		/** \todo Notificar jugadores incluyendo nombre de nuevo jugador **/
	})
	return resp
}

func (cfg *GameConfig) Show() chan string {
	stream := make(chan string)
	go func() {
		if(cfg == nil) {
			stream <- `null`
			close(stream)
			return
		}
		stream <- `{"Rows":` + strconv.Itoa(int(cfg.Rows))
		stream <- `,"Cols":` + strconv.Itoa(int(cfg.Cols))
		stream <- `,"PMin":` + strconv.Itoa(int(cfg.PMin))
		stream <- `,"PMax":` + strconv.Itoa(int(cfg.PMax))
		stream <- `}`
		close(stream)
	} ()
	return stream
}

func (game *Game) Show() chan string {
	stream := make(chan string)
	if(game == nil) {
		go func() {
			stream <- `null`
			close(stream)
		} ()
	} else {
		game.Loop.Async(func(loop *Loop){
			status := game.Status
			players := make([]*Player, 0, len(game.Players))
			copy(players, game.Players)

			// Para el resto no necesita bloquear game
			go func() {

				stream <- `{"status":"` + status.str() + `"`
				stream <- `,"gameId":` + game.Id.str()
				stream <- `,"players":[`
				for i, player := range players {
					if(i>0) {
						stream <- ","
					}
					for chunk := range player.Show() {
						stream <- chunk
					}
				}
				stream <- `]` // array players
	
				stream <- `,"config":`
				for chunk := range game.Config.Show() {
					stream <- chunk
				}
				

				stream <- `,"pieces":[`
				first := true
				for _,piece := range game.Pieces {
					if( first ) {
						first = false
					} else {
						stream <- ","
					}
					for chunk := range piece.Show() {
						stream <- chunk
					}
				}
				// Pieces  map[PieceId]*Piece
				stream <- `]}` // pieces y game
				close(stream)
			} ()
		})
	}
	return stream
}
func (game *Game) Kill(reason string) {
	/** Notificar a games (falta llevar referencia). **/
	/** Notificar a jugadores. **/
	/** Matar piezas. **/
}