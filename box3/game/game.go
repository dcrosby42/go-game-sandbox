package game

import (
	"github.com/dcrosby42/go-game-sandbox/helpers"
	"github.com/go-gl/gl/v3.3-core/gl"
	_ "github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type State struct {
	Width       int
	Height      int
	Shader      ShaderProgram
	Camera      Camera
	Objects     []Object
	Renderables []*helpers.Renderable
	Angle       float32
	Projection  mgl.Mat4
}

type Action struct{}

func Init(s *State) *State {
	if s.Width <= 0 {
		s.Width = 500
	}
	if s.Height <= 0 {
		s.Height = 500
	}

	diffuseShader := helpers.MustMakeDiffuseShader()
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

	// s.Shader = MakeShader()
	// s.Shader.Use()
	s.Projection = mgl.Perspective(mgl.DegToRad(45.0), float32(s.Width)/float32(s.Height), 0.01, 20.0)
	// s.Shader.SetProjection(projection)
	// s.Shader.SetCamera(camera.Matrix())

	s.Camera = Camera{
		Up:    mgl.Vec3{0, 1, 0},
		Eye:   mgl.Vec3{3, 3, 3},
		Focus: mgl.Vec3{0, 0, 0},
	}
	// camera := mgl.LookAtV(mgl.Vec3{3, 3, 3}, mgl.Vec3{0, 0, 0},

	s.Objects = []Object{
		MakeBox(),
		MakeBox(),
	}
	s.Objects[1].Location[0] = -2

	return s
}

func Update(s *State, action *Action) *State {
	// Update box's rotation
	s.Angle += (3.1415926 / 6) / 12
	s.Renderables[0].Rotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})
	s.Renderables[1].LocalRotation = mgl.QuatRotate(s.Angle, mgl.Vec3{0, 1, 0})

	// eye := &s.Camera.Eye
	// eye[1] -= 0.05
	// if eye[1] < 0 {
	// 	eye[1] = 0
	// }
	return s
}

func Draw(s *State) {
	cameraView := s.Camera.Matrix()

	for _, node := range s.Renderables {
		node.Draw(s.Projection, cameraView)
	}
	// s.Shader.Use()
	//
	// s.Shader.SetCamera(s.Camera.Matrix())
	//
	// for _, obj := range s.Objects {
	// 	s.Shader.SetModel(obj.Matrix())
	// 	obj.Draw()
	// }
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
	box.Location = mgl.Vec3{0, 0, 0}
	box.Rotation = mgl.Vec3{0, 0, 0}
	return
}

type Object struct {
	helpers.Drawable
	Location mgl.Vec3
	Rotation mgl.Vec3
}

func (me Object) Matrix() mgl.Mat4 {
	rotX := mgl.HomogRotate3DX(me.Rotation[0])
	rotY := mgl.HomogRotate3DY(me.Rotation[1])
	rotZ := mgl.HomogRotate3DZ(me.Rotation[2])
	rot := rotX.Mul4(rotY).Mul4(rotZ)

	trans := mgl.Translate3D(me.Location[0], me.Location[1], me.Location[2])

	return trans.Mul4(rot)
}

type Camera struct {
	Eye, Focus, Up mgl.Vec3
}

func (me Camera) Matrix() mgl.Mat4 {
	return mgl.LookAtV(
		me.Eye,
		me.Focus,
		me.Up,
	)
}
