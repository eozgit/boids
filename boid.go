package main

import (
	"math"
	"sync"

	"github.com/dhconnelly/rtreego"
)

var (
	tol              = 0.01
	maxVel           = .5
	separationRange  = 7.
	separationWeight = 0.018
	alignmentRange   = 10.
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
		boid.velocity = boid.velocity.Add(velocity)
	}

	boid.velocity = boid.velocity.Limit(maxVel)
}

func (boid *Boid) separation() {
	defer boid.wg.Done()
	velocity := &Vector{}
	centreOfSearchArea := boid.position().Add(boid.velocity.Limit(1).Scale(separationRange * .3))
	neighbours, neighbourCount := getNeighbours(centreOfSearchArea, separationRange, boid.id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			repel := boid.position().Add(neighbour.position().Negate())
			repel = repel.Scale(1 / math.Pow(repel.Magnitude(), 1.5))
			velocity = velocity.Add(repel)
		}
		velocity = velocity.Scale(separationWeight)
	}
	boid.velocityChan <- velocity
}

func (boid *Boid) alignment() {
	defer boid.wg.Done()
	velocity := &Vector{}
	neighbours, neighbourCount := getNeighbours(boid.position(), alignmentRange, boid.id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			velocity = velocity.Add(neighbour.velocity)
		}
		velocity = velocity.Scale(1 / float64(neighbourCount)).Scale(alignmentWeight)
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
			neighbourPositions = neighbourPositions.Add(neighbour.position())
		}
		velocity = neighbourPositions.Scale(1 / float64(neighbourCount)).Add(boid.position().Negate()).Scale(cohesionWeight)
	}
	boid.velocityChan <- velocity
}
