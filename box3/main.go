package main

// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

import (
	"log"
	"runtime"
	"time"

	"github.com/dcrosby42/go-game-sandbox/box3/game"
	"github.com/dcrosby42/go-game-sandbox/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	fps = 60
)

func main() {
	runtime.LockOSThread()
	w := 500
	h := 500

	win, err := window.New(window.Options{
		Title:  "Box",
		Width:  w,
		Height: h,
	})
	if err != nil {
		log.Fatal("Window creation failed: err=%s", err)
		return
	}

	defer glfw.Terminate()

	// INIT
	state := &game.State{Width: w, Height: h}
	state = game.Init(state)
	var action *game.Action

	for !win.ShouldClose() {
		t := time.Now()

		// UPDATE
		// TODO action = &game.Action{}
		state = game.Update(state, action)

		// DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		game.Draw(state)

		win.SwapBuffers()

		// WAIT
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		glfw.PollEvents()
	}
}
