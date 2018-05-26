package window

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Options struct {
	Title     string
	Width     int
	Height    int
	Resizable bool
}

func New(opts Options) (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	if opts.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	win, err := glfw.CreateWindow(opts.Width, opts.Height, opts.Title, nil, nil)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		return nil, err
	}
	glfw.SwapInterval(1) // enable vsync
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	// version := gl.GoStr(gl.GetString(gl.VERSION))
	// log.Println("OpenGL version", version)

	return win, nil
}
