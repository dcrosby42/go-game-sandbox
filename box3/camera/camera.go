package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

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
