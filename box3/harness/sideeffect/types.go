package sideeffect

type Event interface {
	SideEffect() Event
}

type eventBase struct{}

func (me eventBase) SideEffect() Event {
	return me
}

type Error struct {
	eventBase
	Error error
}

func NewError(err error) *Error {
	return &Error{Error: err}
}

type MouseMode_Game struct {
	eventBase
}
type MouseMode_UI struct {
	eventBase
}
