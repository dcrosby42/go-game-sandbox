package main

// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

import (
	"log"
	"runtime"

	"github.com/dcrosby42/go-game-sandbox/box3/harness"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	h, err := harness.New()
	if err != nil {
		log.Fatalf("Harness setup failed. err=%s", err)
		return
	}

	h.Play()
}
