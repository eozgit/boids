package main

import "math/rand"

type Noise struct{}

func (_ *Noise) Delta(boid *Boid) (velocity *Vector) {
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5
	velocity = &Vector{vx, vy}
	velocity = velocity.Scale(noiseWeight)
	return
}
