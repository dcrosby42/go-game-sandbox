package main

// import (
// "log"
// "runtime"
// "github.com/go-gl/gl/v3.3-core/gl"
// _ "github.com/go-gl/gl/v4.1-core/gl"
// "github.com/go-gl/glfw/v3.2/glfw"
// "github.com/nullboundary/glfont"
// )

import (
	"log"
	"runtime"

	// "github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const windowWidth = 1920
const windowHeight = 1080

func init() {
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "glfontExample", glfw.GetPrimaryMonitor(), nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	//load font (fontfile, font scale, window width, window height
	fontFile := "/Library/Fonts/Trebuchet MS.ttf"
	font, err := glfont.LoadFont(fontFile, int32(52), windowWidth, windowHeight, nil)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.5, 0.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//set color and draw text
		font.SetColor(0.0, 1.0, 0.0, 1.0) //r,g,b,a font color
		// font.Printf(100, 100, 1.0, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.") //x,y,scale,string,printf args
		font.Printf(50, 130, 1.0, "ATOMIC OBJECT <-- good ol' Trebuchet!") //x,y,scale,string,printf args

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
