package memotest

type MoveType int8

const (
    Select MoveType = iota
    Deselect
)

type Move struct {
    GameId uint
    From   *Player
    Type   MoveType
}
