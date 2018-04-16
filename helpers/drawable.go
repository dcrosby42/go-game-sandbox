package helpers

import "github.com/go-gl/gl/v3.3-core/gl"

type Drawable struct {
	Mode, Drawable uint32
	First, Count   int32
}

func (me Drawable) Draw() {
	gl.BindVertexArray(me.Drawable)
	gl.DrawArrays(me.Mode, me.First, me.Count)
}
