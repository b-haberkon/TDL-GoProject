package memotest

import (
    //"encoding/json" // Era usando para enviar JSON, ahora se hace por partes con stream
    "fmt"
	"math/rand"
	"strconv"
	//"sync"
	"time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
)

var store		*	session.Store
var playersList	*	Players
var gamesList	*	Games
var rng 		*	rand.Rand

var ej1 []*Symbol = []*Symbol {
    {Text: "A", Pair: 0}, {Text: "a", Pair: 0},
    {Text: "B", Pair: 1}, {Text: "b", Pair: 1},
    {Text: "C", Pair: 2}, {Text: "c", Pair: 2},
    {Text: "D", Pair: 3}, {Text: "d", Pair: 3},
    {Text: "E", Pair: 4}, {Text: "e", Pair: 4},
    {Text: "F", Pair: 5}, {Text: "f", Pair: 5},
    {Text: "G", Pair: 6}, {Text: "g", Pair: 6},
    {Text: "H", Pair: 7}, {Text: "h", Pair: 7},
    {Text: "I", Pair: 8}, {Text: "i", Pair: 8},
    {Text: "J", Pair: 9}, {Text: "j", Pair: 9} }

/******************************************************************************
**** FUNCIONES AUXILIARES PARA LOS CONTROLADORES ******************************
******************************************************************************/

/** @brief Devuelve un canal en el que una gouroutine
 * escribirá un true después del tiempo dado. Pero
 * además es posible reiniciar el tiempo escribiendo
 * en el segundo ganal.
 * @param d Duración del tiempo de espera.
 * \todo ¿Cambiar el canal de reset por un chan time.Duration?
 **/
func timeout(d time.Duration) (chan bool, chan bool) {
	resp := make(chan bool)
	rst  := make(chan bool)
	go func() {
		exit := false
		for {
			select {
			case <- rst:			// Mientras haya algo para leer en reset,
				exit = false		// simplemente se asegura de que exit sea false.
			default:				// Sólo si no hay un reset,
				if(exit) {			// si no hubo un reset entre la última espera y esta iteración,
					resp<-true		// envía que finalizó,
					close(resp)		// cierra el canal,
					return			// y finaliza.
				} else {			// Si es inicial o viene de un reset,
					time.Sleep(d)	// duerme el tiempo necesario (en el que puede llegar un reset),
					exit = true		// y setea que en la próxima iteración (si no hubo un reset), finalice.
				}	
			} // select
		} // for
	} ()
	return resp, rst
}

// Auxiliar para ctrlWrap, es el tipo de la función que será pasada,
// escrita por separado para aumentar la legibilidad.
type ctrlWrapped func(*fiber.Ctx,chan bool) (SStr, []error);

/** @brief Envuelve a una función de controlador que devuelve la respuesta,
 * de forma que si la envuelta no finaliza a tiempo, la envolvente responderá
 * con un error. Además, si la respuesta incluye un error, lo convierte a JSON.
 */
func ctrlWrap(c *fiber.Ctx, ctrl ctrlWrapped ) error {
	to,rst := timeout(10 * time.Second)
	resp   := make(chan SStr)
	go func() { 
		str, errs := ctrl(c,rst)					// Llama al controlador
		if (nil != errs) && (0 != len(errs)) {
			str = showRetErrors(errs)
		}
		resp <- str
		close(resp)
	}()

	chanStr := make(chan SStr)
	go func() {
		done := false
		closed := 0
		for {
			select {
			case _, ok := <- to:
				if(!ok) {
					to = nil
					closed ++
				} else
				if !done {
					done = true
					chanStr <- SStr1(`{"error":{"message":"timeout"}}`)
				}
			case recv, ok := <- resp:
				if( ok) {
					if done {
						fmt.Printf("Out of time answer discarded")
					} else {
						done = true
						chanStr <- recv
					} // done
				} else { // !ok
					resp = nil
					closed ++
				} // ok … !ok
			} // select
			if( closed > 1) {
				return
			}
		} // for
	} () // goroutine
	/**
	select {
	case <- to:
		s = SStr1(`{"error":{"message":"timeout"}}`)
	case x := <- resp:
		s = x
	} */

	/** \todo ¿También poner un timeout en el flat, ya que es lazy? **/
	return c.SendStream( (<- chanStr).fork(NewSStrSpy("SEND",16)).AsIO() ) // s.Flat()
}

// Auxiliar para getPlayerAndId
type PlayerWithId struct{
	Ptr *Player
	Id  PlayerId
}

/** Devuelve el objeto jugador y su id asociado a la sesión, o un error.
 * @param c Contexto de la petición HTTP, de donde busca la información de sesión.
 **/
func getPlayerAndId(c *fiber.Ctx) RetWithError[PlayerWithId] {
	ret  := WithError[PlayerWithId]{}
	resp := NewRetWithError[PlayerWithId]()

	go func() { // La que realmente trae los valores
		var sess *session.Session
		defer func() { resp.SendAndClose(ret) } ()

		// Trae la sesión
		sess, ret.err = store.Get(c);
		if(ret.err != nil) {
			return
		}

		player, ok := (sess.Get("player")).(*Player)	// Trae el valor en la sesión
		ok = ok && ( <- player.IsValid() )				// Debe ser válido
		if(!ok) { // Si no había valor, obtiene uno nuevo
			name := "" /** \todo Obtener o generar anónimoNNN */
			intent := <- playersList.NewPlayer(name, nil)
			if(intent.err != nil) {
				ret.err = intent.err

				return
			}
			player = intent.val
		}

		// Acá ya tengo el player (viejo o nuevo)
		idIntent := <- player.GetId()
		if(idIntent.err != nil) {
			ret.err = idIntent.err
			return
		}
		ret.val.Id  = idIntent.val
		ret.val.Ptr = player
		sess.Set("player",player)
		} () // Fin: La que realmente trae los valores
	return resp
}

// Auxuliar para getGameAndId
type GameWithId struct{
	Ptr *Game
	Id  GameId
}

/** Devuelve el objeto juego y su id asociado a la petición, o un error.
 * @param c Contexto de la petición HTTP, de donde busca la ruta de la petición.
 **/
func getGameAndId(c *fiber.Ctx) (resp RetWithError[GameWithId]) {
	ret  := WithError[GameWithId]{}
	resp =  NewRetWithError[GameWithId]()

	go func() { // La que realmente trae los valores
		defer func() { resp.SendAndClose(ret) } ()

		num, err := strconv.Atoi(c.Params("gameId"))
		if(err != nil) {
			ret.err = err
			return
		}
		ret.val.Id = GameId{num}

		intent := <- gamesList.GetById(ret.val.Id)
		ret.err = intent.err
		ret.val.Ptr = intent.val
	} () // Fin: La que realmente trae los valores
	return resp
}

/******************************************************************************
**** FUNCIONES PÚBLICAS (CONTROLADORES HTTP) **********************************
******************************************************************************/

/** Inicializa lo necesario para el minijuego Memotest. **/
func CtrlStart() {
	rng			= 	rand.New(rand.NewSource(time.Now().UnixNano()))
    store		=	session.New()
	gamesList 	=	NewGames(nil)
	playersList	=	NewPlayers(gamesList,nil)
}

/** @brief Intenta crear un juego. **/
func CreateGame(c *fiber.Ctx) error {
	return ctrlWrap(c, func(c *fiber.Ctx, rstTo chan bool) (SStr, []error) {
		playerInfo := <- getPlayerAndId(c)
		if(playerInfo.err != nil) {
			return nil, []error{playerInfo.err}
		}

		player, playerId := playerInfo.val.Ptr, playerInfo.val.Id
		fmt.Printf("Player %v: creating game…\n", playerId)

		var extra any = nil
		config := GameConfig{
			Rows : uint8(2+rng.Intn(8)),	// 2 a 10
			Cols : uint8(2+rng.Intn(8)),	// 2 a 10
			Syms : ej1,
			PMin : uint8(1+rng.Intn(1)),	// 1 o 2 (simple )
			PMax : uint8(2+rng.Intn(4)) }	// 2 a 6
	
		intent := <- player.NewGame(gamesList, config, extra)
		if( intent.err != nil ) {
			return nil, []error{intent.err}
		} 

		fmt.Printf("Player %v: create game %v (%v×%v)\n",
			playerId, intent.val.Id, config.Cols, config.Rows)
		return intent.val.Show(), nil
	}) // Función interna y llamada
} // Función externa


/** @brief Intenta unirse a un juego. **/
func JoinGame(c *fiber.Ctx) error {
	return ctrlWrap(c, func(c *fiber.Ctx, rstTo chan bool) (SStr, []error) {
		// Lanzo las dos subrutinas
		chanPlayer := getPlayerAndId(c)
		chanGame   := getGameAndId(c)
		// Espero a que completen ambas
		infoPlayer := <- chanPlayer
		infoGame   := <- chanGame

		errs := removeNils[error](infoPlayer.err, infoGame.err)
		if(len(errs)>0) {
			return nil, errs
		}

		player, playerId := infoPlayer.val.Ptr, infoPlayer.val.Id
		game,   gameId   := infoGame.val.Ptr,   infoGame.val.Id

		fmt.Printf("Player %v: joining game %v…\n", playerId, gameId.str() )
		
		intent := <- player.JoinGame(game)
		if(intent.err != nil) {
			return nil, []error{intent.err}
		}
		return intent.val.Show(), nil
	}) // Función interna y llamada
} // Función externa

/** Intenta mostrar un juego. **/
func ShowGame(c *fiber.Ctx) error {
	return ctrlWrap(c, func(c *fiber.Ctx, rstTo chan bool) (SStr, []error) {
		// Lanzo las dos subrutinas
		chanPlayer := getPlayerAndId(c)
		chanGame   := getGameAndId(c)
		// Espero a que completen ambas
		infoPlayer := <- chanPlayer
		infoGame   := <- chanGame
		errs := removeNils[error](infoPlayer.err, infoGame.err)
		if(len(errs)>0) {
			fmt.Printf("DEBUG: ShowGame: %v\n",errs)
			return nil, errs
		}
		fmt.Printf("Player %v: watch game %v…\n",
			infoPlayer.val.Id, infoGame.val.Id.str() )
		game := infoGame.val.Ptr
		data := game.Show()
		return data, nil
	}) // Función interna y llamada
} // Función externa

func SelectPiece(c *fiber.Ctx) error {
    return nil
}

func DeselectPiece(c *fiber.Ctx) error {
    return nil
}

