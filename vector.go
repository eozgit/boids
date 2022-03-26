package main

import "math"

type Vector struct {
	x float64
	y float64
}

func (v *Vector) add(v2 *Vector) *Vector {
	return &Vector{
		x: v.x + v2.x,
		y: v.y + v2.y,
	}
}

func (v *Vector) scale(factor float64) *Vector {
	return &Vector{
		x: v.x * factor,
		y: v.y * factor,
	}
}

func (v *Vector) limit(max float64) *Vector {
	magnitude := v.magnitude()
	if magnitude > max {
		return v.scale(max / magnitude)
	}
	return v
}

func (v *Vector) magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

func (v *Vector) negate() *Vector {
	return &Vector{
		x: -v.x,
		y: -v.y,
	}
}
