package game

var boxPts = []float32{
	// front four points
	-0.5, 0.5, 0.5,
	-0.5, -0.5, 0.5,
	0.5, -0.5, 0.5,
	0.5, 0.5, 0.5,
	// back four points
	-0.5, 0.5, -0.5,
	-0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, 0.5, -0.5,
}

var boxTris = []int{
	// front face
	0, 1, 2,
	0, 2, 3,
	// left face
	4, 5, 1,
	4, 1, 0,
	// back face
	7, 6, 5,
	7, 5, 4,
	// right face
	3, 2, 6,
	3, 6, 7,
	// top face
	4, 0, 3,
	4, 3, 7,
	// bottom face
	1, 5, 6,
	1, 6, 2,
}

func RectPrism(x, y, z, xs, ys, zs float32) []float32 {
	size := len(boxTris) * 3
	prism := make([]float32, size, size)
	for i, pnum := range boxTris {
		prism[i*3] = (xs * boxPts[pnum*3]) + x
		prism[i*3+1] = (ys * boxPts[pnum*3+1]) + y
		prism[i*3+2] = (zs * boxPts[pnum*3+2]) + z
	}
	return prism
}

func Cube(s float32) []float32 {
	return RectPrism(0, 0, 0, s, s, s)
}

func Cube1() []float32 {
	return Cube(1)
}
