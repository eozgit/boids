package main

import "math/rand"

type Noise struct{}

var _ Velocity = (*Noise)(nil)

func (_ *Noise) Delta(boid *Boid) *Vector {
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5
	velocity := &Vector{vx, vy}
	return velocity.Scale(boid.params.noiseWeight.value())
}
