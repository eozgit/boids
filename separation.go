package main

import (
	"math"
)

type Separation struct{}

var _ Velocity = (*Separation)(nil)

func (_ *Separation) Delta(boid *Boid) *Vector {
	velocity := &Vector{}
	centreOfSearchArea := boid.Position().Add(boid.Velocity.Limit(1).Scale(separationRange * .3))
	neighbours, neighbourCount := GetNeighbours(centreOfSearchArea, separationRange, boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	for _, neighbour := range neighbours {
		repel := boid.Position().Add(neighbour.Position().Negate())
		repel = repel.Scale(1 / math.Pow(repel.Magnitude(), 1.5))
		velocity = velocity.Add(repel)
	}

	return velocity.Scale(separationWeight)
}
