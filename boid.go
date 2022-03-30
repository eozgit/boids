package main

import (
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/gravestench/mathlib"
)

const tol = 0.01

type Boid struct {
	Id       int
	position rtreego.Point
	Velocity *mathlib.Vector2
	trail    []mathlib.Vector2
}

func (boid Boid) Bounds() *rtreego.Rect {
	return boid.position.ToRect(tol)
}

func (boid *Boid) Position() *mathlib.Vector2 {
	return mathlib.NewVector2(boid.position[0], boid.position[1])
}

func (boid *Boid) calculateVelocity() {
	var velocityChan = make(chan *mathlib.Vector2, global.velocityComponentCount)

	var wg sync.WaitGroup
	wg.Add(global.velocityComponentCount)
	for _, velocityComponent := range global.velocityComponents {
		component := velocityComponent
		go func() {
			defer wg.Done()
			velocityChan <- component.Delta(boid)
		}()
	}
	wg.Wait()
	close(velocityChan)

	for velocity := range velocityChan {
		boid.Velocity.Add(velocity)
	}

	boid.Velocity.Limit(global.params.maximumVelocity.value())
}

func wrap(position *mathlib.Vector2) {
	switch {
	case position.X < 0:
		position.X += fWidth
	case position.X > fWidth:
		position.X -= fWidth
	}
	switch {
	case position.Y < 0:
		position.Y += fHeight
	case position.Y > fHeight:
		position.Y -= fHeight
	}
}

func (boid *Boid) update(tick int) {
	position := boid.Position()
	boid.trail[tick%global.params.trailLength.value()] = *position

	boid.calculateVelocity()

	position.Add(boid.Velocity)
	wrap(position)
	boid.position = rtreego.Point{position.X, position.Y}
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
		trailPartIndex := i
		go func() {
			defer wg.Done()
			trailPosition := boid.trail[(tick+trailPartIndex)%trailLength]
			x := int(trailPosition.X)
			y := int(trailPosition.Y)
			pixelIndex := (y*Width + x) * 4
			colourValue := byte(255 * float64(trailLength-trailPartIndex) / float64(trailLength))
			trailChan <- newTrailPixel(pixelIndex, colourValue)
		}()
	}
	wg.Wait()
	close(trailChan)
}

func newBoid(id int, position *mathlib.Vector2) *Boid {
	if position == nil {
		px := rand.Float64() * Width
		py := rand.Float64() * Height
		position = mathlib.NewVector2(px, py)
	}
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5

	trailLength := global.params.trailLength.value()
	trail := make([]mathlib.Vector2, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, *position)
	}

	return &Boid{
		Id:       id,
		position: rtreego.Point{position.X, position.Y},
		Velocity: mathlib.NewVector2(vx, vy),
		trail:    trail,
	}
}
