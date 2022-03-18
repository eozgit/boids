package main

import "math"

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Add(v2 *Vector) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector) Scale(factor float64) {
	v.x *= factor
	v.y *= factor
}

func (v *Vector) Limit(max float64) {
	speed := math.Sqrt(square(v.x) + square(v.y))
	if speed > max {
		v.Scale(max / speed)
	}
}

func square(x float64) float64 {
	return math.Pow(x, 2)
}
