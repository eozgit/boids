package main

import (
	"image/color"
	"math"
	"sync"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/num/quat"
)

type Game struct {
	boids []Boid
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	pointChan := make(chan kdbush.Point, boidCount)
	for _, boid := range g.boids {
		wg.Add(1)
		go func(b Boid) {
			defer wg.Done()
			newPosition := b.position.Add(b.velocity)
			b.setPosition(newPosition)
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
		op := &ebiten.DrawImageOptions{}

		oa := boid.velocity.y / boid.velocity.x
		q := quat.Number{Real: oa}
		atan := quat.Atan(q)
		theta := atan.Real
		if boid.velocity.x > 0 {
			theta += math.Pi / 2
		} else {
			theta -= math.Pi / 2
		}
		op.GeoM.Rotate(theta)

		op.GeoM.Translate(boid.position.x, boid.position.y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
