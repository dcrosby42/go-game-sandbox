package main

import (
	"fmt"
	"log"
	// "fmt"

	_ "image/png"

	// _ "net/http/pprof"

	"github.com/dcrosby42/go-game-sandbox/runloop"
	"github.com/dcrosby42/go-game-sandbox/window"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	mainthread.Run(func() {
		// (don't be fooled; this func body is not actually on the main thread, please use mainthread.Call() to do win/gl stuff)
		var (
			err error
			win *glfw.Window
		)

		// Setup the window
		mainthread.Call(func() {
			win, err = window.New(window.Options{
				Width:  800,
				Height: 600,
				Title:  "My GL Win",
			})
		})
		if err != nil {
			log.Fatal(err)
		}

		// MAIN LOOP
		runloop.At60Fps(win.ShouldClose, func(dt float64) (bool, error) {
			fmt.Printf("dt=%v\n", dt)

			return true, nil
		})
		if err != nil {
			log.Fatal(err)
		}
	})
}
