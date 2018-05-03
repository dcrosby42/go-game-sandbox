package main

var rectPts = []float32{
	// front four points
	-0.5, 0.5, 0, // top left
	-0.5, -0.5, 0, // bot left
	0.5, -0.5, 0, // bot right
	0.5, 0.5, 0, // top right
}

var rectTris = []int{
	0, 1, 2,
	0, 2, 3,
}

func Rect(x, y, z, xs, ys, zs float32) []float32 {
	size := len(rectTris) * 3
	expanded := make([]float32, size, size)
	for i, pnum := range rectTris {
		expanded[i*3] = (xs * rectPts[pnum*3]) + x
		expanded[i*3+1] = (ys * rectPts[pnum*3+1]) + y
		expanded[i*3+2] = (zs * rectPts[pnum*3+2]) + z
	}
	return expanded
}

func Square(s float32) []float32 {
	return Rect(0, 0, 0, s, s, s)
}

func Square1() []float32 {
	return Square(1)
}
