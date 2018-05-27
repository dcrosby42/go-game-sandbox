package game

import "github.com/go-gl/glfw/v3.2/glfw"

//go:generate stringer -type=ActionType

type ActionType int

const (
	Tick ActionType = iota
	MouseEnter
	MouseMove
	MouseButton
	MouseScroll
	Keyboard
	Char
	WindowSize
)

type Action struct {
	Type        ActionType
	Tick        *TickAction
	MouseEnter  *MouseEnterAction
	MouseMove   *MouseMoveAction
	MouseButton *MouseButtonAction
	MouseScroll *MouseScrollAction
	Keyboard    *KeyboardAction
	Char        *CharAction
	WindowSize  *WindowSizeAction
}

type TickAction struct {
	Gt float64
	Dt float64
}

type KeyboardAction struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action // Press, Release, Repeat
	Modifier glfw.ModifierKey
}

type CharAction struct {
	Char     rune
	Modifier glfw.ModifierKey
}

func (me CharAction) String() string {
	return string(me.Char)
}

type WindowSizeAction struct {
	FbWidth, FbHeight int
	Width, Height     int
}

type MouseEnterAction struct {
	Entered bool
}

type MouseMoveAction struct {
	PixX, PixY, PixDx, PixDy float32
	X, Y, Dx, Dy             float32
	InBounds                 bool
}

// MouseButton:
//   glfw.MouseButton1
//   glfw.MouseButton2  ...8, Last, Left, Right, Middle
// Action:
//   glfw.Press
//   glfw.Release
// ModifierKey: (bitmask)
//   glfw.ModShift
//   glfw.ModControl
//   glfw.ModAlt
//   glfw.ModSuper (Cmd or Win key)
type MouseButtonAction struct {
	Button   glfw.MouseButton
	Action   glfw.Action
	Modifier glfw.ModifierKey
}

type MouseScrollAction struct {
	X, Y float64
}
