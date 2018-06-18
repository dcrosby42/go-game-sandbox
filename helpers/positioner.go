package helpers

import mgl "github.com/go-gl/mathgl/mgl32"

type Positioner struct {
	// Location positions the object in world space
	Location mgl.Vec3

	// Scale represents how to scale the object when drawing
	Scale mgl.Vec3

	// Rotation is the rotation of the object in world space
	Rotation mgl.Quat

	// LocalRotation is rotation applied to the object in local space
	LocalRotation mgl.Quat

	cached    bool
	transform mgl.Mat4
}

func NewPositioner() *Positioner {
	p := new(Positioner)
	p.cached = false
	p.Scale = mgl.Vec3{1, 1, 1}
	return p
}

// GetTransformMat4 creates a transform matrix: scale * transform
func (me *Positioner) GetTransform() mgl.Mat4 {
	if me.cached {
		return me.transform
	}
	scaleMat := mgl.Scale3D(me.Scale[0], me.Scale[1], me.Scale[2])
	transMat := mgl.Translate3D(me.Location[0], me.Location[1], me.Location[2])
	localRotMat := me.LocalRotation.Mat4()
	rotMat := me.Rotation.Mat4()
	me.transform = rotMat.Mul4(transMat).Mul4(localRotMat).Mul4(scaleMat)
	me.cached = true
	return me.transform
}

func (me *Positioner) SetDirty() {
	me.cached = false
}

func (me *Positioner) IsDirty() bool {
	return !me.cached
}
