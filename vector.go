package main

type Vector struct {
	X, Y float64
}

func (v1 *Vector) add(v2 Vector) {
	v1.X += v2.X
	v1.Y += v2.Y
}

func (v *Vector) scale(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}
