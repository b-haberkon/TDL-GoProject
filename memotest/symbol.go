package memotest

type Symbol struct {
    Text string `json:"text"`	// Texto o emoji que se mostrará (puede ser distinto para piezas del mismo par)
    Pair int    `json:"pair"`	// El par es igual sólo en las piezas que hacen juego
}
