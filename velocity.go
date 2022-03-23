package main

import (
	"fmt"
	"sync"

	"github.com/dhconnelly/rtreego"
)

const (
	maxVel = .6
)

var (
	separationRange  = 10.
	separationWeight = 0.000005
	alignmentRange   = 30.
	alignmentWeight  = .01
	homingWeight     = .016
)

type velocityMethod func()

type VelocityCalculator struct {
	boid         *Boid
	wg           *sync.WaitGroup
	velocityChan chan *Vector
}

func (ops *VelocityCalculator) calculate(boid *Boid) {
	ops.boid = boid
	ops.wg = &sync.WaitGroup{}
	velocityMethods := []velocityMethod{ops.separation, ops.alignment}
	velocityMethodCount := len(velocityMethods)
	ops.velocityChan = make(chan *Vector, velocityMethodCount)

	ops.wg.Add(velocityMethodCount)
	for _, velocityMethod := range velocityMethods {
		velocityMethod()
	}
	ops.wg.Wait()
	close(ops.velocityChan)
	for velocity := range ops.velocityChan {
		ops.boid.velocity.Add(velocity)
	}

	ops.boid.velocity.Limit(maxVel)
}

func (ops *VelocityCalculator) separation() {
	defer ops.wg.Done()
	avgPosNei := &Vector{}
	position := ops.boid.position()
	arr := rt.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(separationRange))
	l := len(arr)
	desc := fmt.Sprintf("sep_%d", l)
	if len(arr) < 2 {
		ops.velocityChan <- &Vector{description: desc}
	} else {
		for _, spa := range arr {
			boid := spa.(Boid)
			if boid.id != ops.boid.id {
				avgPosNei.Add(boid.position())
			}
		}
		avgPosNei.Scale(1 / float64(len(arr)))
		velocity := &Vector{x: position.x - avgPosNei.x, y: position.y - avgPosNei.y, description: desc}
		velocity.Scale(separationWeight)
		ops.velocityChan <- velocity
	}
}

func (ops *VelocityCalculator) alignment() {
	defer ops.wg.Done()
	position := ops.boid.position()
	arr := rt.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(alignmentRange))
	avgVelNei := &Vector{description: "align"}
	for _, spa := range arr {
		boid := spa.(Boid)
		if boid.id != ops.boid.id {
			avgVelNei.Add(boid.velocity)
		}
	}
	avgVelNei.Scale(1 / float64(len(arr)))
	avgVelNei.Scale(alignmentWeight)
	ops.velocityChan <- avgVelNei
}
