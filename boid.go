package main

import (
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
)

const tol = 0.01

type Boid struct {
	position     rtreego.Point
	Id           int
	Velocity     *Vector
	wg           *sync.WaitGroup
	velocityChan chan *Vector
	trail        []Vector
	params       *Parameters
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.position.ToRect(tol)
}

func (boid *Boid) Position() *Vector {
	return &Vector{x: boid.position[0], y: boid.position[1]}
}

func (boid *Boid) calculateVelocity() {
	velocityCalcs := []Velocity{
		&Separation{},
		&Alignment{},
		&Cohesion{},
		&Noise{},
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

	boid.Velocity = boid.Velocity.Limit(boid.params.maxVel.value())
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
	boid.trail[tick%boid.params.trailLength.value()] = *position

	boid.calculateVelocity()

	position = position.Add(boid.Velocity)
	wrap(position)
	boid.position = rtreego.Point{position.x, position.y}
}

type TrailPixel struct {
	pixelIndex  int
	colourValue byte
}

func (boid *Boid) getTrailPixels(tick int, trailChan chan *TrailPixel) {
	var wg sync.WaitGroup
	trailLength := boid.params.trailLength.value()
	wg.Add(trailLength)
	for i := 0; i < trailLength; i++ {
		go func(trailPartIndex int) {
			defer wg.Done()
			trailPosition := boid.trail[(tick+trailPartIndex)%trailLength]
			x := int(trailPosition.x)
			y := int(trailPosition.y)
			pixelIndex := (y*Width + x) * 4
			colourValue := byte(255 * float64(trailLength-trailPartIndex) / float64(trailLength))
			trailChan <- &TrailPixel{pixelIndex, colourValue}
		}(i)
	}
	wg.Wait()
	close(trailChan)
}

func newBoid(id int, params *Parameters) *Boid {
	px := rand.Float64() * Width
	py := rand.Float64() * Height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5

	trailLength := params.trailLength.value()
	trail := make([]Vector, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, Vector{px, py})
	}

	return &Boid{
		Id:       id,
		position: rtreego.Point{px, py},
		Velocity: &Vector{vx, vy},
		trail:    trail,
		params:   params,
	}
}
