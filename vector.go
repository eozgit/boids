package main

import "math"

type Vector struct {
	x           float64
	y           float64
	description string
}

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{
		x: v.x + v2.x,
		y: v.y + v2.y,
	}
}

func (v *Vector) Scale(factor float64) *Vector {
	return &Vector{
		x: v.x * factor,
		y: v.y * factor,
	}
}

func (v *Vector) Limit(max float64) *Vector {
	magnitude := v.Magnitude()
	if magnitude > max {
		return v.Scale(max / magnitude)
	}
	return v
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(square(v.x) + square(v.y))
}

func (v *Vector) Negate() *Vector {
	return &Vector{
		x: -v.x,
		y: -v.y,
	}
}

func square(x float64) float64 {
	return math.Pow(x, 2)
}
