package game

import (
	"math"

	"github.com/dcrosby42/go-game-sandbox/box3/camera"
	"github.com/dcrosby42/go-game-sandbox/helpers"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	Pi                   = math.Pi
	Pi_2                 = math.Pi / 2
	Pi_4                 = math.Pi / 4
	Pi_6                 = math.Pi / 6
	TwoPi                = math.Pi * 2
	cameraMoveSpeed      = 5
	mouseLookSensitivity = 0.25
)

type State struct {
	Width             int
	Height            int
	Camera            camera.Camera
	StartCamera       camera.Camera
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
	cube2.Location = mgl.Vec3{-2, -2, 0}
	cube2.Tex0 = crateTexture

	cube3 := helpers.CreateCube(-0.25, -0.25, -0.25, 0.25, 0.25, 0.25)
	cube3.Shader = diffuseShader
	cube3.Color = mgl.Vec4{0.25, 1.0, 0.25, 1.0}
	cube3.Location = mgl.Vec3{2, 2, 0}
	cube3.Tex0 = crateTexture

	s.Renderables = []*helpers.Renderable{
		cube1,
		cube2,
		cube3,
	}

	s.Projection = mgl.Perspective(Pi_4, float32(s.Width)/float32(s.Height), 0.01, 20.0)

	s.Camera = camera.Camera{
		Position:  mgl.Vec3{0, 0, 7},
		Yaw:       Pi_2,
		Pitch:     0,
		MinPitch:  Pi/-2 + 0.0001,
		MaxPitch:  Pi/2 - 0.0001,
		UseTarget: false,
		Target:    &mgl.Vec3{0, 0, 0},
	}
	s.StartCamera = s.Camera //copy

	s.Camera.Update()

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
		// if movePositionSimple(&s.Camera.Position, &s.CameraMoveControl, speed) {
		if movePositionFps(&s.Camera.Position, &s.Camera.DirFront, &s.Camera.DirLeft, &s.Camera.DirUp, &s.CameraMoveControl, speed) {
			s.Camera.Update()
		}

	case Keyboard:
		updateWasdDirControl(&s.CameraMoveControl, action.Keyboard)

		// Reset Camera
		if action.Keyboard.Key == glfw.Key0 && action.Keyboard.Action == glfw.Press {
			s.Camera = s.StartCamera
			s.Camera.Update()
		}

		if action.Keyboard.Key == glfw.KeyT && action.Keyboard.Action == glfw.Press {
			s.Camera.UseTarget = !s.Camera.UseTarget
			s.Camera.Update()
		}

		// Rotate camera via arrow keys
		if action.Keyboard.Key == glfw.KeyLeft && action.Keyboard.Action == glfw.Press {
			s.Camera.Yaw += Pi_6 / 2
			s.Camera.Update()
		}
		if action.Keyboard.Key == glfw.KeyRight && action.Keyboard.Action == glfw.Press {
			s.Camera.Yaw -= Pi_6 / 2
			s.Camera.Update()
		}
		if action.Keyboard.Key == glfw.KeyUp && action.Keyboard.Action == glfw.Press {
			s.Camera.Pitch += Pi_6 / 2
			s.Camera.Update()
		}
		if action.Keyboard.Key == glfw.KeyDown && action.Keyboard.Action == glfw.Press {
			s.Camera.Pitch -= Pi_6 / 2
			s.Camera.Update()
		}

	case Char:
		// fmt.Printf("game.Update() Char: %s mods=%d\n", action.Char.String(), action.Char.Modifier)
	case MouseEnter:
		// fmt.Printf("game.Update() MouseEnter: %v\n", action.MouseEnter.Entered)
	case MouseMove:
		if action.MouseMove.MouseDrive {
			s.Camera.Yaw -= math.Mod(float64(action.MouseMove.Dx*mouseLookSensitivity), TwoPi)
			s.Camera.Pitch -= float64(action.MouseMove.Dy * mouseLookSensitivity)
			s.Camera.Update()
		}
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
	// cameraView := s.Camera.Matrix
	cameraView := s.Camera.Matrix

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

func movePositionSimple(pos *mgl.Vec3, dirControl *DirControl, dist float32) (changed bool) {
	if dirControl.Up {
		pos[2] -= dist
		changed = true
	}
	if dirControl.Down {
		pos[2] += dist
		changed = true
	}
	if dirControl.Right {
		pos[0] += dist
		changed = true
	}
	if dirControl.Left {
		pos[0] -= dist
		changed = true
	}
	return changed
}

func movePositionFps(pos, front, left, up *mgl.Vec3, dirControl *DirControl, dist float32) (changed bool) {
	if dirControl.Up {
		// pos[2] -= dist
		*pos = pos.Add(front.Mul(dist))
		changed = true
	}
	if dirControl.Down {
		// pos[2] += dist
		*pos = pos.Sub(front.Mul(dist))
		changed = true
	}
	if dirControl.Right {
		// pos[0] += dist
		*pos = pos.Sub(left.Mul(dist))
		changed = true
	}
	if dirControl.Left {
		*pos = pos.Add(left.Mul(dist))
		// pos[0] -= dist
		changed = true
	}
	return changed
}
