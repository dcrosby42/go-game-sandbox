package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position, WorldUp mgl.Vec3
	LookTarget        *mgl.Vec3
}

func (me Camera) Matrix() mgl.Mat4 {
	return mgl.LookAtV(
		me.Position,
		*me.LookTarget,
		me.WorldUp,
	)
}
