package game

import (
	"github.com/dcrosby42/go-game-sandbox/helpers"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type State struct {
	Width   int
	Height  int
	Shader  ShaderProgram
	Camera  Camera
	Objects []Object
}

type Action struct{}

func Init(s *State) *State {
	if s.Width <= 0 {
		s.Width = 500
	}
	if s.Height <= 0 {
		s.Height = 500
	}
	s.Shader = MakeShader()
	s.Shader.Use()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(s.Width)/float32(s.Height), 0.01, 20.0)
	s.Shader.SetProjection(projection)
	// s.Shader.SetCamera(camera.Matrix())

	s.Camera = Camera{
		Up:    mgl32.Vec3{0, 1, 0},
		Eye:   mgl32.Vec3{3, 3, 3},
		Focus: mgl32.Vec3{0, 0, 0},
	}
	// camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0},

	s.Objects = []Object{
		MakeBox(),
		MakeBox(),
	}
	s.Objects[1].Location[0] = -2

	return s
}

func Update(s *State, action *Action) *State {
	// Update box's rotation
	s.Objects[0].Rotation[1] += (3.1415926 / 6) / 12
	s.Objects[1].Rotation[1] += (3.1415926 / 6) / 12
	eye := &s.Camera.Eye
	eye[1] -= 0.05
	if eye[1] < 0 {
		eye[1] = 0
	}
	return s
}

func Draw(s *State) {
	s.Shader.Use()

	s.Shader.SetCamera(s.Camera.Matrix())

	for _, obj := range s.Objects {
		s.Shader.SetModel(obj.Matrix())
		obj.Draw()
	}
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

func MakeBox() (box Object) {
	box.Drawable = helpers.Wireframe{MakeCube()}
	box.Location = mgl32.Vec3{0, 0, 0}
	box.Rotation = mgl32.Vec3{0, 0, 0}
	return
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
