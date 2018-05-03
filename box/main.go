package main

// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

import (
	"log"
	"runtime"
	"time"

	"github.com/dcrosby42/go-game-sandbox/helpers"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 500
	height = 500
	fps    = 30
)

func main() {
	runtime.LockOSThread()

	window := newGlfwWindow()
	defer glfw.Terminate()

	initOpenGL()

	mainView := NewProgram()

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.01, 20.0)

	camera := Camera{
		Up:    mgl32.Vec3{0, 1, 0},
		Eye:   mgl32.Vec3{3, 3, 3},
		Focus: mgl32.Vec3{0, 0, 0},
	}
	// camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0},

	objects := []Object{
		MakeBox(),
		MakeBox(),
	}
	objects[1].Location[0] = -2

	mainView.Use()
	mainView.SetProjection(projection)
	mainView.SetCamera(camera.Matrix())

	for !window.ShouldClose() {
		t := time.Now()

		// Update box's rotation
		objects[0].Rotation[1] += (3.1415926 / 6) / 12
		objects[1].Rotation[1] += (3.1415926 / 6) / 12
		camera.Eye[1] -= 0.05
		if camera.Eye[1] < 0 {
			camera.Eye[1] = 0
		}

		// == BEGIN DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		mainView.Use()

		mainView.SetCamera(camera.Matrix())

		for _, obj := range objects {
			mainView.SetModel(obj.Matrix())
			obj.Draw()
		}

		window.SwapBuffers()
		// == END DRAW

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		glfw.PollEvents()
	}
}

// initializes glfw and returns a Window to use.
func newGlfwWindow() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func MakeCube() helpers.Drawable {
	pts := RectPrism(0, 0, 0, 1, 1, 1)
	tris := int32(len(pts) / 2)
	vao := helpers.MakeVao(pts)
	return &helpers.DrawableVertexArray{
		Mode:     gl.TRIANGLES,
		Drawable: vao,
		First:    0,
		Count:    tris,
	}
}

type Object struct {
	helpers.Drawable
	Location mgl32.Vec3
	Rotation mgl32.Vec3
}

func (me Object) Matrix() mgl32.Mat4 {
	rotX := mgl32.HomogRotate3DX(me.Rotation[0])
	rotY := mgl32.HomogRotate3DY(me.Rotation[1])
	rotZ := mgl32.HomogRotate3DZ(me.Rotation[2])
	rot := rotX.Mul4(rotY).Mul4(rotZ)

	trans := mgl32.Translate3D(me.Location[0], me.Location[1], me.Location[2])

	return trans.Mul4(rot)
}

type Camera struct {
	Eye, Focus, Up mgl32.Vec3
}

func (me Camera) Matrix() mgl32.Mat4 {
	return mgl32.LookAtV(
		me.Eye,
		me.Focus,
		me.Up,
	)
}

func MakeBox() (box Object) {
	// box.Matrix = mgl32.Ident4()
	// box.Drawable = MakeCube()
	box.Drawable = helpers.Wireframe{MakeCube()}
	box.Location = mgl32.Vec3{0, 0, 0}
	box.Rotation = mgl32.Vec3{0, 0, 0}
	return
}
