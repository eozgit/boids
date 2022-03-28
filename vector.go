package main

import "math"

type Vector struct {
	x float64
	y float64
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
	if magnitude <= max {
		return v
	}
	return v.Scale(max / magnitude)
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

func (v *Vector) Negate() *Vector {
	return &Vector{
		x: -v.x,
		y: -v.y,
	}
}
