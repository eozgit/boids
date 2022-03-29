package main

import (
	"math/rand"

	"github.com/gravestench/mathlib"
)

type Noise struct{}

var _ Velocity = (*Noise)(nil)

func (_ *Noise) Delta(boid *Boid) *mathlib.Vector2 {
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5
	velocity := mathlib.NewVector2(vx, vy)
	return velocity.Scale(global.params.noiseWeight.value())
}

func newNoise() *Noise {
	return &Noise{}
}
