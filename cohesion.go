package main

import (
	"github.com/gravestench/mathlib"
)

type Cohesion struct{}

var _ Velocity = (*Cohesion)(nil)

func (_ *Cohesion) Delta(boid *Boid) *mathlib.Vector2 {
	velocity := mathlib.NewVector2(0, 0)
	neighbours, neighbourCount := global.index.GetNeighbours(boid.Position(), global.params.cohesionRange.value(), boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	neighbourPositions := mathlib.NewVector2(0, 0)
	for _, neighbour := range neighbours {
		neighbourPositions.Add(neighbour.Position())
	}

	return neighbourPositions.Scale(1 / float64(neighbourCount)).Subtract(boid.Position()).Scale(global.params.cohesionWeight.value())
}

func newCohesion() *Cohesion {
	return &Cohesion{}
}
