package game

import (
	"fmt"

	"github.com/dcrosby42/go-game-sandbox/box3/camera"
	"github.com/dcrosby42/go-game-sandbox/helpers"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type State struct {
	Width       int
	Height      int
	Camera      camera.Camera
	Renderables []*helpers.Renderable
	Angle       float32
	Projection  mgl.Mat4
}

func Init(s *State) (*State, error) {
	if s.Width <= 0 {
		s.Width = 500
	}
	if s.Height <= 0 {
		s.Height = 500
	}

	diffuseShader, err := helpers.LoadShaderProgramFromFile(
		"shaders/diffuse_texture.vert.glsl",
		"shaders/diffuse_texture.frag.glsl",
	)
	if err != nil {
		return s, err
	}

	crateTexture := helpers.MustLoadTexture("assets/crate1_diffuse.png")

	cube1 := helpers.CreateCube(-0.5, -0.5, -0.5, 0.5, 0.5, 0.5)
	cube1.Shader = diffuseShader
	cube1.Color = mgl.Vec4{1.0, 1.0, 1.0, 1.0}
	cube1.Location = mgl.Vec3{0, 0, 0}
	cube1.Tex0 = crateTexture

	cube2 := helpers.CreateCube(-0.5, -0.5, -0.5, 0.5, 0.5, 0.5)
	cube2.Shader = diffuseShader
	cube2.Color = mgl.Vec4{1.0, 1.0, 1.0, 1.0}
	cube2.Location = mgl.Vec3{-2, 0, 0}
	cube2.Tex0 = crateTexture

	s.Renderables = []*helpers.Renderable{
		cube1,
		cube2,
	}

	s.Projection = mgl.Perspective(mgl.DegToRad(45.0), float32(s.Width)/float32(s.Height), 0.01, 20.0)

	s.Camera = camera.Camera{
		Up:    mgl.Vec3{0, 1, 0},
		Eye:   mgl.Vec3{3, 3, 3},
		Focus: mgl.Vec3{0, 0, 0},
	}

	return s, nil
}

func Update(s *State, action *Action) *State {
	switch action.Type {
	case Tick:
		// Update box's rotation
		s.Angle += (3.1415926 / 6) / 12
		s.Renderables[0].Rotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})
		s.Renderables[1].LocalRotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})

		// descend camera
		eye := &s.Camera.Eye
		eye[1] -= 0.05
		if eye[1] < 0 {
			eye[1] = 0
		}
	case MouseMove:
		if action.MouseMove.InBounds {
			fmt.Printf("MouseMove(%f,%f, %v)\n", action.MouseMove.X, action.MouseMove.Y, action.MouseMove.InBounds)
		}
	}

	return s
}

func Draw(s *State) {
	cameraView := s.Camera.Matrix()

	for _, node := range s.Renderables {
		node.Draw(s.Projection, cameraView)
	}
}
