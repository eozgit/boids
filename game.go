package main

import (
	"image/color"
	"sync"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	homingWeight    = .01
	alignmentWeight = .002
	alignmentRadius = 30.
	op              = &ebiten.DrawImageOptions{}
)

type Game struct {
	boids []*Boid
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	pointChan := make(chan kdbush.Point, boidCount)
	for _, boid := range g.boids {
		wg.Add(1)
		go func(b *Boid) {
			defer wg.Done()
			vHoming := getHomingVelocity(b.position)
			vHoming.Scale(homingWeight)
			vAlignment := getAlignmentVelocity(b)
			vAlignment.Scale(alignmentWeight)
			newVelocity := b.velocity.Add(vHoming).Add(vAlignment)
			newVelocity.Limit(.5)
			b.velocity.x = newVelocity.x
			b.velocity.y = newVelocity.y
			newPosition := b.position.Add(b.velocity)
			b.setPosition(newPosition)
			b.calculateAngle()
			pointChan <- &kdbush.SimplePoint{X: b.position.x, Y: b.position.y}
		}(boid)
	}
	wg.Wait()
	close(pointChan)
	points = []kdbush.Point{}
	for point := range pointChan {
		points = append(points, point)
	}
	bush = kdbush.NewBush(points, boidCount)
	return nil
}

func getAlignmentVelocity(boid *Boid) *Vector {
	arr := bush.Within(&kdbush.SimplePoint{X: boid.position.x, Y: boid.position.y}, alignmentRadius)
	v := &Vector{}
	for _, i := range arr {
		v = v.Add(boids[i].velocity)
	}
	l := float64(len(arr))
	v.x = v.x / l
	v.y = v.y / l
	return v
}

func getHomingVelocity(position *Vector) *Vector {
	x, y := 0., 0.
	fWidth, fHeight := float64(width), float64(height)
	outOfBoundsLeft := position.x < 0
	outOfBoundsRight := position.x > fWidth
	outOfBoundsTop := position.y < 0
	outOfBoundsBottom := position.y > fHeight
	if outOfBoundsLeft {
		x = -position.x / fWidth
	} else if outOfBoundsRight {
		x = -(position.x - fWidth) / fWidth
	}
	if outOfBoundsTop {
		y = -position.y / fHeight
	} else if outOfBoundsBottom {
		y = -(position.y - fHeight) / fHeight
	}
	return &Vector{x, y}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	for _, boid := range g.boids {
		op.GeoM.Reset()

		op.GeoM.Rotate(boid.angle)

		op.GeoM.Translate(boid.position.x, boid.position.y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
