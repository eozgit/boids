package main

import (
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
)

var (
	tol    = 0.01
	maxVel = .5
)

type Boid struct {
	rtreego.Point
	Id           int
	Velocity     *Vector
	wg           *sync.WaitGroup
	velocityChan chan *Vector
	trail        []Vector
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.ToRect(tol)
}

func (boid *Boid) Position() *Vector {
	return &Vector{x: boid.Point[0], y: boid.Point[1]}
}

func (boid *Boid) calculateVelocity() {
	velocityCalcs := []Velocity{
		&Separation{},
		&Alignment{},
		&Cohesion{},
	}
	velocityCalcCount := len(velocityCalcs)
	boid.velocityChan = make(chan *Vector, velocityCalcCount)

	boid.wg = &sync.WaitGroup{}
	boid.wg.Add(velocityCalcCount)
	for _, velocityCalc := range velocityCalcs {
		go func(velCalc Velocity) {
			defer boid.wg.Done()
			boid.velocityChan <- velCalc.Delta(boid)
		}(velocityCalc)
	}
	boid.wg.Wait()
	close(boid.velocityChan)
	for velocity := range boid.velocityChan {
		boid.Velocity = boid.Velocity.Add(velocity)
	}

	boid.Velocity = boid.Velocity.Limit(maxVel)
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
	position := boid.Position()
	boid.trail[tick%trailLength] = *position

	boid.calculateVelocity()

	position = position.Add(boid.Velocity)
	wrap(position)
	boid.Point = rtreego.Point{position.x, position.y}
}

func newBoid(id int) *Boid {
	px := rand.Float64() * Width
	py := rand.Float64() * Height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5

	trail := make([]Vector, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, Vector{px, py})
	}

	return &Boid{
		Id:       id,
		Point:    rtreego.Point{px, py},
		Velocity: &Vector{vx, vy},
		trail:    trail,
	}
}
