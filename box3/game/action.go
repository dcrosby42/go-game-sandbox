package game

//go:generate stringer -type=ActionType

type ActionType int

const (
	Tick ActionType = iota
	MouseMove
)

type Action struct {
	Type      ActionType
	Tick      *TickAction
	MouseMove *MouseMoveAction
}

type TickAction struct {
	Gt float64
	Dt float64
}

type MouseMoveAction struct {
	PixX, PixY, PixDx, PixDy float32
	X, Y, Dx, Dy             float32
	InBounds                 bool
}
