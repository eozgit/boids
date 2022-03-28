package main

type Velocity interface {
	Delta(boid *Boid) *Vector
}
