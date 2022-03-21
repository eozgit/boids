package main

import (
	"errors"
	"image/color"
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	op = &ebiten.DrawImageOptions{}
)

type Game struct {
	boids []*Boid
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	pointChan := make(chan rtreego.Spatial, boidCount)
	for _, boid := range g.boids {
		wg.Add(1)
		go func(b *Boid) {
			defer wg.Done()

			velocityCalc := VelocityCalculator{}
			velocityCalc.calculate(b)

			b.position.Add(b.velocity)
			b.calculateAngle()
			pointChan <- Point{rtreego.Point{b.position.x, b.position.y}, b}
		}(boid)
	}
	wg.Wait()
	close(pointChan)
	points = []rtreego.Spatial{}
	for point := range pointChan {
		points = append(points, point)
	}
	rt = rtreego.NewTree(2, 5, 500, points...)
	return nil
}

func (g *Game) getBoidById(id int) (*Boid, error) {
	for _, boid := range g.boids {
		if boid.id == id {
			return boid, nil
		}
	}
	return nil, errors.New("Boid not found.")
}

func setRGB(matrix *ebiten.ColorM, red int, green int, blue int) {
	// Reset RGB (not Alpha) 0 forcibly
	matrix.Scale(0, 0, 0, 1)

	// Set color
	r := float64(255) / 0xff
	g := float64(0) / 0xff
	b := float64(0) / 0xff
	matrix.Translate(r, g, b, 0)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	b0, _ := g.getBoidById(0)
	re := rtreego.Point{b0.position.x, b0.position.y}.ToRect(separationRange)
	arr := rt.SearchIntersect(re)

	for _, boid := range g.boids {
		op.GeoM.Reset()
		op.ColorM.Reset()

		op.GeoM.Rotate(boid.angle)

		op.GeoM.Translate(boid.position.x, boid.position.y)

		if boid.id == 0 {
			setRGB(&op.ColorM, 255, 0, 0)
		} else {
			for _, spa := range arr {
				point := spa.(Point)
				if boid.id == point.boid.id {
					setRGB(&op.ColorM, 0, 255, 255)
					break
				}
			}
		}

		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
