package main

import (
	"math"
)

type Boid struct {
	id       int
	strId    string
	position *Vector
	velocity *Vector
	angle    float64
}

func (boid *Boid) calculateAngle() {
	boid.angle = math.Atan2(boid.velocity.y, boid.velocity.x) + math.Pi/2
}
