package harness

import (
	"time"

	"github.com/dcrosby42/go-game-sandbox/box3/game"
	"github.com/dcrosby42/go-game-sandbox/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/go-gl/glow/gl"
	// "github.com/go-gl/glow/gl"
)

type Harness struct {
	fps           int
	width, height int
	win           *glfw.Window
	state         *game.State
	lastGameTime  float64
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
	return har, nil
}

func (me *Harness) Play() {
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
		action := &game.Action{
			Type: game.Tick,
			Tick: &game.TickAction{Gt: gameTime, Dt: dt},
		}
		me.state = game.Update(me.state, action)

		// DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		game.Draw(me.state)

		me.win.SwapBuffers()

		// WAIT
		time.Sleep(time.Second/time.Duration(me.fps) - time.Since(t))
	}
}
