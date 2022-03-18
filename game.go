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

			velocityCalc := VelocityCalculator{}
			velocityCalc.calculate(b)

			b.position.Add(b.velocity)
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
