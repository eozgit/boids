package main

type Velocity interface {
	Delta(boid *Boid) *Vector
}

var (
	separationRange  = 11.
	separationWeight = 0.017
	alignmentRange   = 19.
	alignmentWeight  = .02
	cohesionRange    = 19.
	cohesionWeight   = .001
)
