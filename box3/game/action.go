package game

//go:generate stringer -type=ActionType

type ActionType int

const (
	Tick ActionType = iota
)

type Action struct {
	Type ActionType
	Tick *TickAction
}

type TickAction struct {
	Gt float64
	Dt float64
}
