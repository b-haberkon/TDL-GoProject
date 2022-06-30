package memotest

type Symbol struct {
    Text string `json:"text"`	// Texto o emoji que se mostrará (puede ser distinto para piezas del mismo par)
    Pair int    `json:"pair"`	// El par es igual sólo en las piezas que hacen juego
}

func shuffleSymbols(org []*Symbol, amount int) []*Symbol {
	/** \todo Agrupar de a pares y otras cosas. **/
	dest := make([]*Symbol, amount)
	l := len(org)
	copy(dest, org[0:min(amount,l)])
	return dest
}

