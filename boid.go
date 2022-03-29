package main

import (
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
)

const tol = 0.01

type Boid struct {
	Id       int
	position rtreego.Point
	Velocity *Vector
	trail    []Vector
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.position.ToRect(tol)
}

func (boid *Boid) Position() *Vector {
	return &Vector{x: boid.position[0], y: boid.position[1]}
}

func (boid *Boid) calculateVelocity() {
	var velocityChan = make(chan *Vector, global.velocityComponentCount)

	var wg sync.WaitGroup
	wg.Add(global.velocityComponentCount)
	for _, velocityComponent := range global.velocityComponents {
		go func(velComp Velocity) {
			defer wg.Done()
			velocityChan <- velComp.Delta(boid)
		}(velocityComponent)
	}
	wg.Wait()
	close(velocityChan)

	for velocity := range velocityChan {
		boid.Velocity = boid.Velocity.Add(velocity)
	}

	boid.Velocity = boid.Velocity.Limit(global.params.maxVel.value())
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
	boid.trail[tick%global.params.trailLength.value()] = *position

	boid.calculateVelocity()

	position = position.Add(boid.Velocity)
	wrap(position)
	boid.position = rtreego.Point{position.x, position.y}
}

type TrailPixel struct {
	pixelIndex  int
	colourValue byte
}

func newTrailPixel(pixelIndex int, colourValue byte) *TrailPixel {
	return &TrailPixel{pixelIndex, colourValue}
}

func (boid *Boid) getTrailPixels(tick int, trailChan chan *TrailPixel) {
	var wg sync.WaitGroup
	trailLength := global.params.trailLength.value()
	wg.Add(trailLength)
	for i := 0; i < trailLength; i++ {
		go func(trailPartIndex int) {
			defer wg.Done()
			trailPosition := boid.trail[(tick+trailPartIndex)%trailLength]
			x := int(trailPosition.x)
			y := int(trailPosition.y)
			pixelIndex := (y*Width + x) * 4
			colourValue := byte(255 * float64(trailLength-trailPartIndex) / float64(trailLength))
			trailChan <- newTrailPixel(pixelIndex, colourValue)
		}(i)
	}
	wg.Wait()
	close(trailChan)
}

func newBoid(id int) *Boid {
	px := rand.Float64() * Width
	py := rand.Float64() * Height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5

	trailLength := global.params.trailLength.value()
	trail := make([]Vector, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, Vector{px, py})
	}

	return &Boid{
		Id:       id,
		position: rtreego.Point{px, py},
		Velocity: &Vector{vx, vy},
		trail:    trail,
	}
}
