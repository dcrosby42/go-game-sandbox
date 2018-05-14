package camera

import (
	"fmt"
	"math"

	"github.com/dcrosby42/go-game-sandbox/helpers"
	mgl "github.com/go-gl/mathgl/mgl32"
)

var worldUp = mgl.Vec3{0, 1, 0}

type Camera struct {
	Matrix mgl.Mat4

	Position                 mgl.Vec3
	MinPitch, MaxPitch       float64
	Pitch, Yaw               float64
	DirFront, DirLeft, DirUp mgl.Vec3
	UseTarget                bool
	Target                   *mgl.Vec3

	DebugUpdates bool
}

var v3s = helpers.Vec3String

func (me *Camera) Update() {
	if me.UseTarget {
		// Set Matrix based on an Position + arbitrary world location at Target
		me.Matrix = mgl.LookAtV(
			me.Position,
			*me.Target,
			worldUp,
		)
		fmt.Printf("Camera.Update() pos=%s target=%s\n", v3s(&me.Position), v3s(me.Target))

	} else {
		// constrain pitch
		me.Pitch = float64(mgl.Clamp(float32(me.Pitch), float32(me.MinPitch), float32(me.MaxPitch)))

		//
		// Set Matrix based on Position + Pitch and Yaw
		//
		// "front" represents a forward-pointing vector based on pitch and yaw
		front := mgl.Vec3{
			float32(math.Cos(me.Pitch) * math.Cos(me.Yaw)),
			float32(math.Sin(me.Pitch)),
			float32(math.Cos(me.Pitch) * -math.Sin(me.Yaw)), // -sin makes it so positive yaw pivots around y-axis from x+ (right) toward z- (into screen).  Left hand coord system
		}.Normalize()
		// "left" is perpendicular to "front" in the local plane of the camera
		left := mgl.Vec3{}.Sub(front.Cross(worldUp).Normalize())
		// "up" vector is perpendicular to the local plane
		up := front.Cross(left).Normalize()
		// look at our own nose:
		target := me.Position.Add(front)
		me.Matrix = mgl.LookAtV(
			me.Position,
			target,
			up,
		)

		// keep track of our current vecs in case someone on the outside is interested
		me.DirFront = front
		me.DirLeft = left
		me.DirUp = up

		if me.DebugUpdates {
			fmt.Printf("Camera.Update() pos=%s yaw=%.2f pitch=%.2f front=%s left=%s up=%s\n", v3s(&me.Position), me.Yaw, me.Pitch, v3s(&me.DirFront), v3s(&me.DirLeft), v3s(&me.DirUp))
		}
	}
}
