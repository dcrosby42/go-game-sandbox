package main

// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

import (
	"fmt"
	"log"
	"runtime"
	"strings"
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

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	box := MakeCube()

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.01, 20.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	var angle float32

	for !window.ShouldClose() {
		t := time.Now()

		// Update the model-level transform matrix:
		angle += (3.1415926 / 6) / 12
		model = mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0})

		// BEGIND DRAW
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		// Set the model xform into the shader:
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		// Draw box
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		box.Draw()
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

		glfw.PollEvents()
		window.SwapBuffers()
		// END DRAW

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
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
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

const (
	vertexShaderSource = `
    #version 410

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

    in vec3 vert;

    void main() {
        gl_Position = projection * camera * model * vec4(vert, 1);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
)

// var vertexShader = `
// #version 330
//
// uniform mat4 projection;
// uniform mat4 camera;
// uniform mat4 model;
//
// in vec3 vert;
// in vec2 vertTexCoord;
//
// out vec2 fragTexCoord;
//
// void main() {
//     fragTexCoord = vertTexCoord;
//     gl_Position = projection * camera * model * vec4(vert, 1);
// }
// ` + "\x00"

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
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
