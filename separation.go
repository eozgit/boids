package main

import (
	"math"
)

type Separation struct{}

func (_ *Separation) Delta(boid *Boid) (velocity *Vector) {
	velocity = &Vector{}
	centreOfSearchArea := boid.Position().Add(boid.Velocity.Limit(1).Scale(separationRange * .3))
	neighbours, neighbourCount := GetNeighbours(centreOfSearchArea, separationRange, boid.Id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			repel := boid.Position().Add(neighbour.Position().Negate())
			repel = repel.Scale(1 / math.Pow(repel.Magnitude(), 1.5))
			velocity = velocity.Add(repel)
		}
		velocity = velocity.Scale(separationWeight)
	}
	return
}
