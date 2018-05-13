package game

import (
	"github.com/dcrosby42/go-game-sandbox/box3/camera"
	"github.com/dcrosby42/go-game-sandbox/helpers"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const cameraMoveSpeed = 5

type State struct {
	Width             int
	Height            int
	Camera            camera.Camera
	Renderables       []*helpers.Renderable
	Angle             float32
	Projection        mgl.Mat4
	CameraMoveControl DirControl
}
type DirControl struct {
	Up, Left, Down, Right bool
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
		WorldUp: mgl.Vec3{0, 1, 0},
		// Eye:   mgl.Vec3{3, 3, 3},
		Position:   mgl.Vec3{3, 0, 3},
		LookTarget: &mgl.Vec3{0, 0, 0},
	}

	return s, nil
}

func Update(s *State, action *Action) *State {
	switch action.Type {
	case Tick:
		// Update box's rotation
		// s.Angle += (3.1415926 / 6) / 12
		// s.Renderables[0].Rotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})
		// s.Renderables[1].LocalRotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})

		// descend camera
		// eye := &s.Camera.Eye
		// eye[1] -= 0.05
		// if eye[1] < 0 {
		// 	eye[1] = 0
		// }
		// move camera
		speed := float32(cameraMoveSpeed * action.Tick.Dt)
		updatePosition(&s.Camera.Position, &s.CameraMoveControl, speed)

	case Keyboard:
		updateWasdDirControl(&s.CameraMoveControl, action.Keyboard)

		if action.Keyboard.Key == glfw.Key0 && action.Keyboard.Action == glfw.Press {
			s.Camera.Position = mgl.Vec3{0, 0, 0}
		}

	case Char:
		// fmt.Printf("game.Update() Char: %s mods=%d\n", action.Char.String(), action.Char.Modifier)
	case MouseEnter:
		// fmt.Printf("game.Update() MouseEnter: %v\n", action.MouseEnter.Entered)
	case MouseMove:
		// if action.MouseMove.InBounds {
		// fmt.Printf("MouseMove(%f,%f, %v)\n", action.MouseMove.X, action.MouseMove.Y, action.MouseMove.InBounds)
		// }
	case MouseButton:
		// fmt.Printf("game.Update() MouseButton: %#v\n", action.MouseButton)
	case MouseScroll:
		// fmt.Printf("game.Update() MouseScroll: %#v\n", action.MouseScroll)

	}

	return s
}

func Draw(s *State) {
	cameraView := s.Camera.Matrix()

	for _, node := range s.Renderables {
		node.Draw(s.Projection, cameraView)
	}
}

func updateWasdDirControl(wasd *DirControl, ka *KeyboardAction) {
	pressed := false
	switch ka.Action {
	case glfw.Press:
		pressed = true
	case glfw.Release:
		pressed = false
	default:
		return
	}

	switch ka.Key {
	case glfw.KeyW:
		wasd.Up = pressed
	case glfw.KeyA:
		wasd.Left = pressed
	case glfw.KeyS:
		wasd.Down = pressed
	case glfw.KeyD:
		wasd.Right = pressed
	}
}
func updateArrowDirControl(wasd *DirControl, ka *KeyboardAction) {
	pressed := false
	switch ka.Action {
	case glfw.Press:
		pressed = true
	case glfw.Release:
		pressed = false
	default:
		return
	}

	switch ka.Key {
	case glfw.KeyUp:
		wasd.Up = pressed
	case glfw.KeyLeft:
		wasd.Left = pressed
	case glfw.KeyDown:
		wasd.Down = pressed
	case glfw.KeyRight:
		wasd.Right = pressed
	}
}

func updatePosition(pos *mgl.Vec3, dirControl *DirControl, dist float32) {
	if dirControl.Up {
		pos[2] -= dist
	}
	if dirControl.Down {
		pos[2] += dist
	}
	if dirControl.Right {
		pos[0] += dist
	}
	if dirControl.Left {
		pos[0] -= dist
	}
}
