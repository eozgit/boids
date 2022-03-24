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
	separationWeight = 0.008
	alignmentRange   = 10.
	alignmentWeight  = .03
)

type velocityMethod func()

type Boid struct {
	rtreego.Point
	id           int
	velocity     *Vector
	angle        float64
	wg           *sync.WaitGroup
	velocityChan chan *Vector
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

func (boid *Boid) calculateVelocity() {
	velocityMethods := []velocityMethod{
		boid.separation,
		boid.alignment,
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
	centreOfSearchArea := boid.position().Add(boid.velocity.Limit(1).Scale(separationRange * .7))
	neighbours := search(centreOfSearchArea, separationRange)
	if len(neighbours) < 2 {
		boid.velocityChan <- velocity
	} else {
		for _, abstract := range neighbours {
			neighbour := abstract.(Boid)
			if neighbour.id == boid.id {
				continue
			}
			repel := boid.position().Add(neighbour.position().Negate())
			repel = repel.Scale(1 / math.Pow(repel.Magnitude(), 1.5))
			velocity = velocity.Add(repel)
		}
		velocity = velocity.Scale(separationWeight)
		boid.velocityChan <- velocity
	}
}

func (boid *Boid) alignment() {
	defer boid.wg.Done()
	neighbours := search(boid.position(), alignmentRange)
	velocity := &Vector{}
	for _, abstract := range neighbours {
		neighbour := abstract.(Boid)
		if boid.id == neighbour.id {
			continue
		}
		velocity = velocity.Add(neighbour.velocity)
	}
	velocity = velocity.Scale(1 / float64(len(neighbours))).Scale(alignmentWeight)
	boid.velocityChan <- velocity
}
