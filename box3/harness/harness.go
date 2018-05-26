package harness

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dcrosby42/go-game-sandbox/box3/game"
	"github.com/dcrosby42/go-game-sandbox/box3/harness/sideeffect"
	"github.com/dcrosby42/go-game-sandbox/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Harness struct {
	fps                 int
	winWidth, winHeight int
	fbWidth, fbHeight   int
	win                 *glfw.Window
	state               *game.State
	lastGameTime        float64
	cursor              CursorState

	DebugInput       bool
	DebugSideEffects bool
}

func New() (*Harness, error) {
	winWidth := 500
	winHeight := 500

	win, err := window.New(window.Options{
		Title:     "Box",
		Width:     winWidth,
		Height:    winHeight,
		Resizable: true,
	})
	if err != nil {
		return nil, err
	}

	fbWidth, fbHeight := win.GetFramebufferSize()

	har := &Harness{
		fps:          120,
		winWidth:     winWidth,
		winHeight:    winHeight,
		fbWidth:      fbWidth,
		fbHeight:     fbHeight,
		win:          win,
		state:        nil,
		lastGameTime: 0,
	}

	har.DebugInput = false
	har.DebugSideEffects = true

	// GLFW input guid: http://www.glfw.org/docs/latest/input_guide.html
	win.SetKeyCallback(har.KeyCallback)
	win.SetCharModsCallback(har.CharModsCallback)
	win.SetMouseButtonCallback(har.MouseButtonCallback)
	win.SetCursorPosCallback(har.CursorPosCallback)
	win.SetCursorEnterCallback(har.CursorEnterCallback)
	win.SetScrollCallback(har.ScrollCallback)
	win.SetFramebufferSizeCallback(har.FramebufferSizeCallback)

	state := &game.State{Width: fbWidth, Height: fbHeight}
	state, sideEffect := game.Init(state)
	har.state = state

	err = har.HandleSideEffect(sideEffect)
	if err != nil {
		return nil, err
	}

	// har.MouseModeGame()

	return har, nil
}

func (me *Harness) Play() {
	action := game.Action{
		Type: game.Tick,
		Tick: &game.TickAction{},
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !me.win.ShouldClose() {
		t := time.Now()
		gameTime := glfw.GetTime()

		glfw.PollEvents()

		dt := gameTime - me.lastGameTime
		if dt > 0.2 {
			dt = 0.2
		}
		me.lastGameTime = gameTime

		// UPDATE
		action.Tick.Gt = gameTime
		action.Tick.Dt = dt
		me.ApplyUpdate(&action)

		// DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		game.Draw(me.state)

		me.win.SwapBuffers()

		// WAIT
		time.Sleep(time.Second/time.Duration(me.fps) - time.Since(t))
	}
}

func (me *Harness) ApplyUpdate(action *game.Action) {
	var e sideeffect.Event
	me.state, e = game.Update(me.state, action)
	err := me.HandleSideEffect(e)
	if err != nil {
		panic(fmt.Sprintf("FAILED ApplyUpdate() err=%s", err)) // TODO c'mon man, we can do better than panic
	}
}

func (me *Harness) HandleSideEffect(e sideeffect.Event) error {
	if e == nil {
		return nil
	}
	if me.DebugSideEffects {
		fmt.Printf("Harness.HandleSideEffect(): %v\n", reflect.TypeOf(e))
	}
	switch event := e.(type) {
	case *sideeffect.Error:
		return event.Error
	case *sideeffect.MouseMode_Game:
		me.MouseModeGame()
	case *sideeffect.MouseMode_UI:
		me.MouseModeUI()
	default:

	}

	return nil

}

func (me *Harness) MouseModeGame() {
	me.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (me *Harness) MouseModeUI() {
	me.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

func (me *Harness) ScrollCallback(w *glfw.Window, xoff, yoff float64) {
	action := game.Action{
		Type: game.MouseScroll,
		MouseScroll: &game.MouseScrollAction{
			X: xoff,
			Y: yoff,
		},
	}
	me.ApplyUpdate(&action)
	if me.DebugInput {
		fmt.Printf("Harness.ScrollCallback() xoff=%f yoff=%f\n", xoff, yoff)
	}
}

func (me *Harness) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, maction glfw.Action, mod glfw.ModifierKey) {
	action := game.Action{
		Type: game.MouseButton,
		MouseButton: &game.MouseButtonAction{
			Button:   button,
			Action:   maction,
			Modifier: mod,
		},
	}

	me.ApplyUpdate(&action)

	if me.DebugInput {
		fmt.Printf("Harness.MouseCallback() button=%d action=%d mod=%d\n", button, maction, mod)
	}

}

func (me *Harness) CursorPosCallback(w *glfw.Window, x float64, y float64) {
	xpos := float32(x)
	ypos := float32(y)
	dx := xpos - me.cursor.x
	dy := ypos - me.cursor.y
	if !me.cursor.tracking {
		dx = 0
		dy = 0
		me.cursor.tracking = true
	}
	me.cursor.x = xpos
	me.cursor.y = ypos

	w2 := float32(me.winWidth) / 2
	h2 := float32(me.winHeight) / 2
	nx := (xpos - w2) / w2
	ny := -(ypos - h2) / h2
	inbounds := nx >= -1.0 && nx <= 1.0 && ny <= 1.0 && ny >= -1.0
	ndx := float32(dx) / w2
	ndy := float32(dy) / h2
	action := game.Action{
		Type: game.MouseMove,
		MouseMove: &game.MouseMoveAction{
			PixX:     xpos,
			PixY:     ypos,
			PixDx:    dx,
			PixDy:    dy,
			X:        nx, // mgl.Clamp(nx, -1, 1),
			Y:        ny, //mgl.Clamp(ny, -1, 1),
			Dx:       ndx,
			Dy:       ndy,
			InBounds: inbounds,
		},
	}

	me.ApplyUpdate(&action)
	if me.DebugInput {
		fmt.Printf("Harness.CursorPosCallback(): Pix(%.2f, %.2f) PixD(%.2f, %.2f) Norm(%.4f, %.4f) NormD(%.4f, %.4f) InBounds=%v\n", xpos, ypos, dx, dy, nx, ny, ndx, ndy, inbounds)
	}
}

func (me *Harness) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if me.DebugInput {
		fmt.Printf("Harness.KeyCallback() key='%s' (%d) scancode=%d action=%d mods=%d\n", glfw.GetKeyName(key, scancode), key, scancode, action, mods)
	}

	kaction := game.Action{
		Type: game.Keyboard,
		Keyboard: &game.KeyboardAction{
			Key:      key,
			Scancode: scancode,
			Action:   action,
			Modifier: mods,
		},
	}
	me.ApplyUpdate(&kaction)

	//XXX
	// if key == glfw.KeyEscape && action == glfw.Press {
	// 	if me.firstPersonMouse {
	// 		me.firstPersonMouse = false
	// 		me.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	// 	} else {
	// 		me.firstPersonMouse = true
	// 		me.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	// 	}
	// }
}

func (me *Harness) CharModsCallback(w *glfw.Window, char rune, mods glfw.ModifierKey) {
	chaction := game.Action{
		Type: game.Char,
		Char: &game.CharAction{
			Char:     char,
			Modifier: mods,
		},
	}
	me.ApplyUpdate(&chaction)

	if me.DebugInput {
		fmt.Printf("Harness.CharModsCallback() char='%s'\n", string(char))
	}
}

func (me *Harness) CursorEnterCallback(w *glfw.Window, entered bool) {
	eaction := game.Action{
		Type: game.MouseEnter,
		MouseEnter: &game.MouseEnterAction{
			Entered: entered,
		},
	}
	me.ApplyUpdate(&eaction)

	if me.DebugInput {
		fmt.Printf("Harness.CursorEnterCallback() entered=%v\n", entered)
	}

	// if entered {
	// 	me.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	// } else {
	// 	me.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	// }

}

type CursorState struct {
	x, y     float32
	tracking bool
}

func (me *Harness) FramebufferSizeCallback(w *glfw.Window, fbWidth, fbHeight int) {
	me.fbWidth = fbWidth
	me.fbHeight = fbHeight
	me.winWidth, me.winHeight = me.win.GetSize()
	waction := game.Action{
		Type: game.WindowSize,
		WindowSize: &game.WindowSizeAction{
			Width:    me.winWidth,
			Height:   me.winHeight,
			FbWidth:  fbWidth,
			FbHeight: fbHeight,
		},
	}
	me.ApplyUpdate(&waction)

	if me.DebugInput {
		fmt.Printf("Harness.FramebufferSizeCallback(%d, %d)\n", fbWidth, fbHeight)
	}
}
