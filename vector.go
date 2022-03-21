package main

import "math"

type Vector struct {
	x           float64
	y           float64
	description string
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
	magnitude := v.Magnitude()
	if magnitude > max {
		v.Scale(max / magnitude)
	}
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(square(v.x) + square(v.y))
}

func square(x float64) float64 {
	return math.Pow(x, 2)
}
