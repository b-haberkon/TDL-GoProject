package memotest

type Symbol struct {
    Text string `json:"text"`	// Texto o emoji que se mostrará (puede ser distinto para piezas del mismo par)
    Pair uint   `json:"pair"`	// El par es igual sólo en las piezas que hacen juego
}

func selectSymbols(org [][]*Symbol, amount int) []*Symbol {
	/** \todo Si se incorpora dificultad, aquí debería seleccionarlos. **/
	// Hacer una copia de los pares
	nPairs := len(org)
	pairs := make([][]*Symbol, nPairs)
	copy(pairs, org)
	// Mezclar (Fisher-Yates) los pares
	for i:=0; i<nPairs; i++ {
		j        := rng.Intn(nPairs)
		temp     := pairs[j]
		pairs[j]  = pairs[i]
		pairs[i]  = temp
	}
	// Copiar los primeros amount/2 pares a un slice unidimensional de fichas
	pieces := make([]*Symbol, 0, amount)
	for i:=0; i<(amount/2); i++ {
		for _,piece := range pairs[i] {
			pieces = append(pieces, piece)
		}
	}
	// Mezclar (Fisher-Yates) las fichas
	for i:=0; i<amount; i++ {
		j         := rng.Intn(amount)
		temp      := pieces[j]
		pieces[j]  = pieces[i]
		pieces[i]  = temp
	}
	return pieces
}

