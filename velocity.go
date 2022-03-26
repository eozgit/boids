package main

type Velocity interface {
	Delta(boid *Boid) *Vector
}

var (
	separationRange  = 6.
	separationWeight = .02
	alignmentRange   = 19.
	alignmentWeight  = .01
	cohesionRange    = 19.
	cohesionWeight   = .0004
	noiseWeight      = .03
)
