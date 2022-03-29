package main

import (
	"math"

	"github.com/gravestench/mathlib"
)

type Separation struct{}

var _ Velocity = (*Separation)(nil)

func (_ *Separation) Delta(boid *Boid) *mathlib.Vector2 {
	velocity := mathlib.NewVector2(0, 0)
	separationRange := global.params.separationRange.value()
	centreOfSearchArea := boid.Position().Add(boid.Velocity.Limit(1).Scale(separationRange * .3))
	neighbours, neighbourCount := global.index.GetNeighbours(centreOfSearchArea, separationRange, boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	for _, neighbour := range neighbours {
		repel := boid.Position().Subtract(neighbour.Position())
		repel = repel.Scale(1 / math.Pow(repel.Length(), 1.5))
		velocity.Add(repel)
	}

	return velocity.Scale(global.params.separationWeight.value())
}

func newSeparation() *Separation {
	return &Separation{}
}
