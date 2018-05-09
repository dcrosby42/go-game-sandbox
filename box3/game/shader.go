package game

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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

type ShaderProgram struct {
	program        uint32
	uni_projection int32
	uni_camera     int32
	uni_model      int32
}

func MakeShader() (prog ShaderProgram) {
	p := NewGlProgram()
	prog.program = p
	prog.uni_projection = gl.GetUniformLocation(p, gl.Str("projection\x00"))
	prog.uni_camera = gl.GetUniformLocation(p, gl.Str("camera\x00"))
	prog.uni_model = gl.GetUniformLocation(p, gl.Str("model\x00"))
	return
}

func (me *ShaderProgram) Use() {
	gl.UseProgram(me.program)
}
func (me *ShaderProgram) SetProjection(proj mgl32.Mat4) {
	gl.UniformMatrix4fv(me.uni_projection, 1, false, &proj[0])
}
func (me *ShaderProgram) SetCamera(camera mgl32.Mat4) {
	gl.UniformMatrix4fv(me.uni_camera, 1, false, &camera[0])
}
func (me *ShaderProgram) SetModel(model mgl32.Mat4) {
	gl.UniformMatrix4fv(me.uni_model, 1, false, &model[0])
}

// gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

func NewGlProgram() uint32 {
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
