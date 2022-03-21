package main

import (
	"fmt"
	"sync"

	"github.com/dhconnelly/rtreego"
)

const (
	maxVel                 = .6
	velocityComponentCount = 3
)

var (
	separationRange  = 4.
	separationWeight = 0.00005
	alignmentRange   = 30.
	alignmentWeight  = .01
	homingWeight     = .016
)

type VelocityCalculator struct {
	boid         *Boid
	wg           *sync.WaitGroup
	velocityChan chan *Vector
}

func (ops *VelocityCalculator) calculate(boid *Boid) {
	ops.boid = boid
	ops.wg = &sync.WaitGroup{}
	ops.velocityChan = make(chan *Vector, velocityComponentCount)

	ops.wg.Add(velocityComponentCount)
	go ops.separation()
	go ops.alignment()
	go ops.homing()
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
	arr := rt.SearchIntersect(rtreego.Point{ops.boid.position.x, ops.boid.position.y}.ToRect(separationRange))
	l := len(arr)
	desc := fmt.Sprintf("sep_%d", l)
	if len(arr) < 2 {
		ops.velocityChan <- &Vector{description: desc}
	} else {
		for _, spa := range arr {
			point := spa.(Point)
			if point.boid.id != ops.boid.id {
				avgPosNei.Add(point.boid.position)
			}
		}
		avgPosNei.Scale(1 / float64(len(arr)))
		velocity := &Vector{x: ops.boid.position.x - avgPosNei.x, y: ops.boid.position.y - avgPosNei.y, description: desc}
		velocity.Scale(separationWeight)
		ops.velocityChan <- velocity
	}
}

func (ops *VelocityCalculator) alignment() {
	defer ops.wg.Done()
	arr := rt.SearchIntersect(rtreego.Point{ops.boid.position.x, ops.boid.position.y}.ToRect(alignmentRange))
	avgVelNei := &Vector{description: "align"}
	for _, spa := range arr {
		point := spa.(Point)
		if point.boid.id != ops.boid.id {
			avgVelNei.Add(point.boid.velocity)
		}
	}
	avgVelNei.Scale(1 / float64(len(arr)))
	avgVelNei.Scale(alignmentWeight)
	ops.velocityChan <- avgVelNei
}

func (ops *VelocityCalculator) homing() {
	defer ops.wg.Done()
	vector := &Vector{homingComponent(ops.boid.position.x, ops.boid.velocity.x, width), homingComponent(ops.boid.position.y, ops.boid.velocity.y, height), "hom"}
	vector.Scale(homingWeight)
	ops.velocityChan <- vector
}

func homingComponent(pos float64, vel float64, dim int) (velocity float64) {
	fDim := float64(dim)
	if pos < 0 {
		velocity = -pos / fDim
		if vel > 0 {
			velocity /= vel
		}
	} else if pos > fDim {
		velocity = (fDim - pos) / fDim
		if vel < 0 {
			velocity /= -vel
		}
	}
	return
}
