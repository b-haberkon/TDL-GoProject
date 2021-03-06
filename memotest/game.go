package memotest

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
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
	Syms [][]*Symbol
	PMin uint8
	PMax uint8
}

type Game struct {
	Loop	      *Loop
	Id		      GameId
	Config	      GameConfig
	Extra	      any
	Players	      []*Player
	Pieces        map[PieceId]*Piece
	Points 	      map[*Player]int
	Status        GameStatus
	inGamePieces  int
}

func NewGame(id GameId, config GameConfig, extra any) *Game {
	game 		:= 	Game{
		Loop	:	nil,
		Id		:	id,
		Config	:	config,
		Extra	:	extra,
		Players	:	make([]*Player, 0, 2),
		Pieces	:	make(map[PieceId]*Piece),
		Points  :   make(map[*Player]int),
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
		// Asumimos que son pares, pieza central no existe
		nPieces := int(game.Config.Rows) * int(game.Config.Cols)
		centralRow := game.Config.Rows
		centralCol := game.Config.Cols
		if( (nPieces % 2) > 0) {
			nPieces--
			// Si son impares, remueve la central para que queden pares:
			centralRow = game.Config.Rows/2;
			centralCol = game.Config.Cols/2;
			}
		game.inGamePieces = nPieces
		symbols := selectSymbols(game.Config.Syms, nPieces)
		var row, col uint8
		pieceId := PieceId{0}
		pos := 0
		for row=0; row<game.Config.Rows; row++ {
			for col=0; col<game.Config.Cols; col++ {
				pieceId.inc()
				isCentral := (centralRow == row) && (centralCol == col);
				var piece *Piece
				if(isCentral) {
					piece = NewEmptyPiece(pieceId, row, col)
				} else {
					piece = NewPiece(pieceId, row, col, symbols[pos], &game)
					pos++
				}
				game.Pieces[pieceId] = piece
			} // for col
		} // for row
	})
	return &game
}

/** @brief Ejecuta la función fn en el bucle del juego (si no es null),
 * o en una goroutine que devuelve un error (si es null).
 * \todo Refactorizar con similares, usando interfaces. **/
func GameAsync[T any](game *Game,resp RetWithError[T], fn LoopFn) RetWithError[T] {
	if (game == nil) {
		go func() {
			ret := WithError[T]{}
			ret.err = errors.New("Null game")
			resp.SendAndClose( ret )
		} ()
	} else {
		game.Loop.Async(fn)
	}
	return resp
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
	return AsyncVal(game != nil)
}

func (game *Game) Join(player *Player) RetWithError[*Game] {
	resp := NewRetWithError[*Game]()
	GameAsync(game, resp, func(loop *Loop) {
		if(player == nil || game.Status == Dead) {
			resp.SendNewAndClose(nil, errors.New("Invalid player or game") )
			return
		}
		/** \todo game.IsValid() player.IsValdi() **/
		/** \todo game not end **/

		if(game.Status == Ended) {
			// Si ya terminó, lo devuelve para que pueda verlo
			// pero no lo suma
			resp.SendNewAndClose(game, errors.New("Game ended"))
			return
		}
		game.Points[player] = 0
		var amount uint8 = uint8( len(game.Players) )
		if(amount >= game.Config.PMax) {
			txt := fmt.Sprintf("Game is full, there are %v, max: %v",
				amount, game.Config.PMax)
			resp.SendNewAndClose(nil, errors.New(txt))
			return
		}

		game.Players = append(game.Players, player)
		resp.SendNewAndClose(game, nil)
		fmt.Printf("Game %v join player %v\n",game.Id,player.Id)
		/** \todo Notificar jugadores incluyendo nombre de nuevo jugador **/

		// Si estaba esperando llega al mínimo, y justo llega, empezar.
		if(game.Config.PMin == amount + 1) && (Waiting == game.Status) {
			/** \todo Notificar, etc.*/
			game.Status = Playing
		}

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


func (game *Game) Show(player *Player, playerId PlayerId) chan string {
	stream := make(chan string)
	if(game == nil) {
		go func() {
			stream <- `null`
			close(stream)
		} ()
	} else {
		game.Loop.Async(func(loop *Loop){
			status := game.Status
			players := make([]*Player, len(game.Players)) // Extranañamente, debe haber n elementos para que copy funcione
			copy(players, game.Players)
			points := make(map[*Player]int)
			for player,point := range game.Points {
				points[player] = point
			}

			// Para el resto no necesita bloquear game
			go func() {

				stream <- `{"status":"` + status.str() + `"`
				stream <- `,"gameId":` + game.Id.str()
				stream <- `,"playerId":` + playerId.str()
				stream <- `,"players":[`
				for i, player := range players {
					if(i>0) {
						stream <- ","
					}
					pointsJson := `"Points":` + strconv.Itoa(points[player])
					for chunk := range player.ShowWith(pointsJson) {
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

func (game *Game) getPiece(pieceId PieceId) RetWithError[*Piece] {
	resp := NewRetWithError[*Piece]()
	GameAsync(game, resp, func(loop *Loop) {
		ret := WithError[*Piece]{nil,nil}
		ret.val = game.Pieces[pieceId]
		if( ret.val == nil) {
			ret.err = errors.New("The piece does not exists")
		}
		resp.SendAndClose(ret);
	})
	return resp
}

func (game *Game) PieceRemoved(piece *Piece, player *Player, points int) {
	game.Loop.Async(func(loop *Loop) {
		game.inGamePieces --
		game.Points[player] += points
		if(game.inGamePieces < 1) {
			// No me intersa la respuesta, pero debo leer o se bloquea game.end()
			// Debo leerla asincrónicamente, o nunca entraré en game.end
			// al no poder salir de acá hasta leer.
			go func() { <- game.end() } ()
		}
	})
}

func (game *Game) end() RetWithError[bool] {
	resp := NewRetWithError[bool]()
	GameAsync(game, resp, func(loop *Loop) {
		var wg sync.WaitGroup
		// Pasar a End a todos los jugadores primero, para que no envíen nada a las piezas.
		for _, player := range game.Players {
			wg.Add(1)
			player.End(&wg)	// No me importa esperar, y la función no devuelve nada.
		}
		wg.Wait() // Espero a que todos hayan terminado para matar definitivamente a las piezas.
		//game.Players = []*Player{}
		for _, piece := range game.Pieces {
			piece.Die()	// No me importa esperar, y la función no devuelve nada.
		}
		game.Pieces = make(map[PieceId]*Piece)
		game.Status = Ended
		resp.SendNewAndClose(true,nil)
	})
	return resp
}