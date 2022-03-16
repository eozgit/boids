package main

import (
	"github.com/MadAppGang/kdbush"
)

type Boid struct {
	id       int
	strId    string
	position *Vector
	velocity *Vector
}

func (boid *Boid) setPosition(position *Vector) {
	boid.position.x = position.x
	boid.position.y = position.y
	points = append(points, &kdbush.SimplePoint{X: position.x, Y: position.y})
}
