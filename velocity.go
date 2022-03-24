package main

import (
	"fmt"
	"math"
	"sync"
)

const (
	maxVel = .5
)

var (
	separationRange  = 7.
	separationWeight = 0.008
	alignmentRange   = 10.
	alignmentWeight  = .03
)

type velocityMethod func()

type VelocityCalculator struct {
	boid         *Boid
	wg           *sync.WaitGroup
	velocityChan chan *Vector
}

func (calc *VelocityCalculator) calculate(boid *Boid) {
	calc.boid = boid
	calc.wg = &sync.WaitGroup{}
	velocityMethods := []velocityMethod{
		calc.separation,
		calc.alignment,
	}
	velocityMethodCount := len(velocityMethods)
	calc.velocityChan = make(chan *Vector, velocityMethodCount)

	calc.wg.Add(velocityMethodCount)
	for _, velocityMethod := range velocityMethods {
		velocityMethod()
	}
	calc.wg.Wait()
	close(calc.velocityChan)
	for velocity := range calc.velocityChan {
		calc.boid.velocity = calc.boid.velocity.Add(velocity)
	}

	calc.boid.velocity = calc.boid.velocity.Limit(maxVel)
}

func (calc *VelocityCalculator) separation() {
	defer calc.wg.Done()
	desc := fmt.Sprintf("sep_%d", boidCount)
	deltaVelocity := &Vector{description: desc}
	centreOfSearchArea := calc.boid.position().Add(calc.boid.velocity.Limit(1).Scale(separationRange / 2))
	boidsWithinRange := search(centreOfSearchArea, separationRange)
	boidCount := len(boidsWithinRange)
	if boidCount < 2 {
		calc.velocityChan <- deltaVelocity
	} else {
		for _, abstract := range boidsWithinRange {
			potentialNeighbour := abstract.(Boid)
			if potentialNeighbour.id != calc.boid.id {
				velocity := calc.boid.position().Add(potentialNeighbour.position().Negate())
				velocity = velocity.Scale(1 / math.Pow(velocity.Magnitude(), 1.5))
				deltaVelocity = deltaVelocity.Add(velocity)
			}
		}
		deltaVelocity = deltaVelocity.Scale(separationWeight)
		calc.velocityChan <- deltaVelocity
	}
}

func (calc *VelocityCalculator) alignment() {
	defer calc.wg.Done()
	position := calc.boid.position()
	arr := search(position, alignmentRange)
	avgVelNei := &Vector{description: "align"}
	for _, spa := range arr {
		boid := spa.(Boid)
		if boid.id != calc.boid.id {
			avgVelNei = avgVelNei.Add(boid.velocity)
		}
	}
	avgVelNei = avgVelNei.Scale(1 / float64(len(arr))).Scale(alignmentWeight)
	calc.velocityChan <- avgVelNei
}
