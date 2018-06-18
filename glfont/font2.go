package glfont

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

// A Font allows rendering of text to an OpenGL context.
type Font2 struct {
	*Font
}

//LoadFont loads the specified font at the given scale.
func LoadFont2(file string, scale int32, windowWidth int, windowHeight int, shaderCompiler ShaderCompilerFunc) (*Font2, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// Configure the default font vertex and fragment shaders
	if shaderCompiler == nil {
		shaderCompiler = newProgram
	}
	program, err := shaderCompiler(vertexFont2Shader, fragmentFontShader)
	if err != nil {
		panic(err)
	}

	// Activate corresponding render state
	gl.UseProgram(program)

	//set screen resolution
	resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
	gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))

	font, err := LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
	if err != nil {
		return nil, err
	}

	return &Font2{font}, nil
}

//Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font2) Tprintf(x, y float32, scale float32, transmat mgl.Mat4, fs string, argv ...interface{}) error {

	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	lowChar := rune(32)

	//setup blending mode
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Activate corresponding render state
	gl.UseProgram(f.program)
	//set text color
	gl.Uniform4f(gl.GetUniformLocation(f.program, gl.Str("textColor\x00")), f.color.r, f.color.g, f.color.b, f.color.a)
	//set screen resolution
	//resUniform := gl.GetUniformLocation(f.program, gl.Str("resolution\x00"))
	//gl.Uniform2f(resUniform, float32(2560), float32(1440))

	// transform matrix
	gl.UniformMatrix4fv(gl.GetUniformLocation(f.program, gl.Str("transmat\x00")), 1, false, &transmat[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(f.vao)

	// Iterate through all characters in string
	for i := range indices {

		//get rune
		runeIndex := indices[i]

		//skip runes that are not in font chacter range
		if int(runeIndex)-int(lowChar) > len(f.fontChar) || runeIndex < lowChar {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		//find rune in fontChar list
		ch := f.fontChar[runeIndex-lowChar]

		//calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*scale
		ypos := y - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width) * scale
		h := float32(ch.height) * scale

		//set quad positions
		var x1 = xpos
		var x2 = xpos + w
		var y1 = ypos
		var y2 = ypos + h

		//setup quad array
		var vertices = []float32{
			//  X, Y, Z, U, V
			// Front
			x1, y1, 0.0, 0.0,
			x2, y1, 1.0, 0.0,
			x1, y2, 0.0, 1.0,
			x1, y2, 0.0, 1.0,
			x2, y1, 1.0, 0.0,
			x2, y2, 1.0, 1.0}

		// Render glyph texture over quad
		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		// Update content of VBO memory
		gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

		//BufferSubData(target Enum, offset int, data []byte)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
		// Render quad
		gl.DrawArrays(gl.TRIANGLES, 0, 24)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	//clear opengl textures and programs
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)

	return nil
}
