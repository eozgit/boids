package main

import "github.com/gravestench/mathlib"

type Velocity interface {
	Delta(boid *Boid) *mathlib.Vector2
}
