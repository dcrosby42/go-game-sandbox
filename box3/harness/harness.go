package harness

import (
	"time"

	"github.com/dcrosby42/go-game-sandbox/box3/game"
	"github.com/dcrosby42/go-game-sandbox/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Harness struct {
	fps           int
	width, height int
	win           *glfw.Window
	state         *game.State
	lastGameTime  float64
	mouse         Mouse
}

func New() (*Harness, error) {
	w := 500
	h := 500

	win, err := window.New(window.Options{
		Title:  "Box",
		Width:  w,
		Height: h,
	})
	if err != nil {
		return nil, err
	}

	state := &game.State{Width: w, Height: h}
	state, err = game.Init(state)
	if err != nil {
		return nil, err
	}

	har := &Harness{
		fps:          60,
		width:        w,
		height:       h,
		win:          win,
		state:        state,
		lastGameTime: 0,
	}

	win.SetCursorPosCallback(har.CursorPosCallback)

	return har, nil
}

func (me *Harness) Play() {
	action := game.Action{
		Type: game.Tick,
		Tick: &game.TickAction{},
	}
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
		me.state = game.Update(me.state, &action)

		// DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		game.Draw(me.state)

		me.win.SwapBuffers()

		// WAIT
		time.Sleep(time.Second/time.Duration(me.fps) - time.Since(t))
	}
}

func (me *Harness) CursorPosCallback(w *glfw.Window, x float64, y float64) {
	xpos := float32(x)
	ypos := float32(y)
	dx := xpos - me.mouse.x
	dy := ypos - me.mouse.y
	if !me.mouse.tracking {
		dx = 0
		dy = 0
		me.mouse.tracking = true
	}
	me.mouse.x = xpos
	me.mouse.y = ypos

	w2 := float32(me.width) / 2
	h2 := float32(me.height) / 2
	nx := (xpos - w2) / w2
	ny := -(ypos - h2) / h2
	inbounds := nx >= -1.0 && nx <= 1.0 && ny <= 1.0 && ny >= -1.0
	ndx := float32(dx) / w2
	ndy := float32(dy) / h2
	action := game.Action{
		Type: game.MouseMove,
		MouseMove: &game.MouseMoveAction{
			PixX: xpos, PixY: ypos, PixDx: dx, PixDy: dy,
			X: mgl.Clamp(nx, -1, 1), Y: mgl.Clamp(ny, -1, 1), Dx: ndx, Dy: ndy,
			InBounds: inbounds,
		},
	}
	me.state = game.Update(me.state, &action)
}

type Mouse struct {
	x, y     float32
	tracking bool
}
