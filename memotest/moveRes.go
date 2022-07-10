package memotest

type MoveResultStatus uint8
type MoveResult struct {
	Status MoveResultStatus
	Pieces []*Piece
}

const (
    Selection 		MoveResultStatus = iota
    Match
	Unmatch
	Inexistent
	Blocked
)

var MoveResultTexts = map[MoveResultStatus]string {
    Selection     : "Selection",       // La pieza fue seleccionada
    Match         : "Match",           // Las piezas fueron emparejadas
	Unmatch       : "Unmatch",         // Las piezas no coincidían
	Inexistent    : "Inexistent",      // La pieza no existe
	Blocked       : "Blocked",         // La pieza está bloqueada
}

func NewMoveResult(initialState MoveResultStatus) *MoveResult {
	return &MoveResult{initialState, make([]*Piece,0,2)}
}

func (status MoveResultStatus) str() string { return MoveResultTexts[status]; }

func (res MoveResult) Show() chan string {
	stream := make(chan string)
	go func() {
		stream <- `{"state":"` + res.Status.str() + `"`
		stream <- `,"pieces":[`
		for i, piece := range res.Pieces {
			if(i>0) {
				stream <- ","
			}
			for chunk := range piece.Show() {
				stream <- chunk
			}
		}
		stream <- `]}`
		close(stream)
	} ()
	return stream;
}
