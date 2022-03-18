package main

import "math"

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{v.x + v2.x, v.y + v2.y}
}

func (v *Vector) Scale(factor float64) {
	v.x = v.x * factor
	v.y = v.y * factor
}

func (v *Vector) Limit(max float64) {
	speed := math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
	if speed > max {
		factor := max / speed
		v.x = v.x * factor
		v.y = v.y * factor
	}
}
