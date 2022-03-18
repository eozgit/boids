package main

import (
	"sync"

	"github.com/MadAppGang/kdbush"
)

const velocityComponentCount = 2

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
	go ops.homing()
	go ops.alignment()
	ops.wg.Wait()
	close(ops.velocityChan)
	for velocity := range ops.velocityChan {
		ops.boid.velocity.Add(velocity)
	}

	ops.boid.velocity.Limit(.5)
}

func (ops *VelocityCalculator) homing() {
	defer ops.wg.Done()
	vector := &Vector{homingComponent(ops.boid.position.x, width), homingComponent(ops.boid.position.y, height)}
	vector.Scale(homingWeight)
	ops.velocityChan <- vector
}

func homingComponent(pos float64, dim int) float64 {
	fDim := float64(dim)
	if pos < 0 {
		return -pos / fDim
	} else if pos > fDim {
		return -(pos - fDim) / fDim
	}
	return 0
}

func (ops *VelocityCalculator) alignment() {
	defer ops.wg.Done()
	arr := bush.Within(&kdbush.SimplePoint{X: ops.boid.position.x, Y: ops.boid.position.y}, alignmentRadius)
	vector := &Vector{}
	for _, i := range arr {
		neighbour := boids[i]
		if neighbour.id != ops.boid.id {
			vector.Add(neighbour.velocity)
		}
	}
	vector.Scale(1 / float64(len(arr)))
	vector.Scale(alignmentWeight)
	ops.velocityChan <- vector
}
