package main

type Boid struct {
	id       int
	strId    string
	position *Vector
	velocity *Vector
}

func (boid *Boid) setPosition(position *Vector) {
	boid.position.x = position.x
	boid.position.y = position.y
}
