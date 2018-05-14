package helpers

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func Vec3String(v *mgl.Vec3) string {
	return fmt.Sprintf("[%.4f %.4f %.4f]", v[0], v[1], v[2])
}
