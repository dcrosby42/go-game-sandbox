package helpers

import "github.com/go-gl/gl/v3.3-core/gl"

type Drawable interface {
	Draw()
}

type DrawableVertexArray struct {
	Mode, Drawable uint32
	First, Count   int32
}

func (me DrawableVertexArray) Draw() {
	gl.BindVertexArray(me.Drawable)
	gl.DrawArrays(me.Mode, me.First, me.Count)
}

type Wireframe struct {
	Child Drawable
}

func (me Wireframe) Draw() {
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	me.Child.Draw()
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}
