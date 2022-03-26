package main

import (
	"math"
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
)

var (
	tol              = 0.01
	maxVel           = .5
	separationRange  = 7.
	separationWeight = 0.018
	alignmentRange   = 13.
	alignmentWeight  = .043
	cohesionRange    = 9.
	cohesionWeight   = .00006
)

type velocityMethod func()

type Boid struct {
	rtreego.Point
	id           int
	velocity     *Vector
	wg           *sync.WaitGroup
	velocityChan chan *Vector
	trail        []Vector
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.ToRect(tol)
}

func (boid *Boid) position() *Vector {
	return &Vector{x: boid.Point[0], y: boid.Point[1]}
}

func (boid *Boid) calculateVelocity() {
	velocityMethods := []velocityMethod{
		boid.separation,
		boid.alignment,
		boid.cohesion,
	}
	velocityMethodCount := len(velocityMethods)
	boid.velocityChan = make(chan *Vector, velocityMethodCount)

	boid.wg = &sync.WaitGroup{}
	boid.wg.Add(velocityMethodCount)
	for _, velocityMethod := range velocityMethods {
		velocityMethod()
	}
	boid.wg.Wait()
	close(boid.velocityChan)
	for velocity := range boid.velocityChan {
		boid.velocity = boid.velocity.add(velocity)
	}

	boid.velocity = boid.velocity.limit(maxVel)
}

func (boid *Boid) separation() {
	defer boid.wg.Done()
	velocity := &Vector{}
	centreOfSearchArea := boid.position().add(boid.velocity.limit(1).scale(separationRange * .3))
	neighbours, neighbourCount := getNeighbours(centreOfSearchArea, separationRange, boid.id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			repel := boid.position().add(neighbour.position().negate())
			repel = repel.scale(1 / math.Pow(repel.magnitude(), 1.5))
			velocity = velocity.add(repel)
		}
		velocity = velocity.scale(separationWeight)
	}
	boid.velocityChan <- velocity
}

func (boid *Boid) alignment() {
	defer boid.wg.Done()
	velocity := &Vector{}
	neighbours, neighbourCount := getNeighbours(boid.position(), alignmentRange, boid.id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			velocity = velocity.add(neighbour.velocity)
		}
		velocity = velocity.scale(1 / float64(neighbourCount)).scale(alignmentWeight)
	}
	boid.velocityChan <- velocity
}

func (boid *Boid) cohesion() {
	defer boid.wg.Done()
	velocity := &Vector{}
	neighbours, neighbourCount := getNeighbours(boid.position(), cohesionRange, boid.id)
	if neighbourCount > 0 {
		neighbourPositions := &Vector{}
		for _, neighbour := range neighbours {
			neighbourPositions = neighbourPositions.add(neighbour.position())
		}
		velocity = neighbourPositions.scale(1 / float64(neighbourCount)).add(boid.position().negate()).scale(cohesionWeight)
	}
	boid.velocityChan <- velocity
}

func wrap(position *Vector) {
	switch {
	case position.x < 0:
		position.x += fWidth
	case position.x > fWidth:
		position.x -= fWidth
	}
	switch {
	case position.y < 0:
		position.y += fHeight
	case position.y > fHeight:
		position.y -= fHeight
	}
}

func (boid *Boid) update(tick int) {
	position := boid.position()
	boid.trail[tick%trailLength] = *position

	boid.calculateVelocity()

	position = position.add(boid.velocity)
	wrap(position)
	boid.Point = rtreego.Point{position.x, position.y}
}

func newBoid(id int) *Boid {
	px := rand.Float64() * width
	py := rand.Float64() * height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5

	trail := make([]Vector, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, Vector{px, py})
	}

	return &Boid{
		id:       id,
		Point:    rtreego.Point{px, py},
		velocity: &Vector{vx, vy},
		trail:    trail,
	}
}
