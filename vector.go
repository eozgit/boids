package main

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{v.x + v2.x, v.y + v2.y}
}
