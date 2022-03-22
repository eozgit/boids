package main

import (
	"math"

	"github.com/dhconnelly/rtreego"
)

var tol = 0.01

type Boid struct {
	rtreego.Point
	id       int
	velocity *Vector
	angle    float64
}

func (boid *Boid) calculateAngle() {
	boid.angle = math.Atan2(boid.velocity.y, boid.velocity.x) + math.Pi/2
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.ToRect(tol)
}

func (boid *Boid) position() *Vector {
	return &Vector{x: boid.Point[0], y: boid.Point[1]}
}
