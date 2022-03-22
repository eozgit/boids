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
	rt    *rtreego.Rtree
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	boidChan := make(chan *Boid, boidCount)
	for _, boid := range g.boids {
		wg.Add(1)
		go func(b *Boid) {
			defer wg.Done()

			velocityCalc := VelocityCalculator{}
			velocityCalc.calculate(b)

			position := b.position()
			position.Add(b.velocity)
			b.Point = rtreego.Point{position.x, position.y}
			b.calculateAngle()
			boidChan <- b
		}(boid)
	}
	wg.Wait()
	close(boidChan)

	boids := []*Boid{}
	points := []rtreego.Spatial{}
	for boid := range boidChan {
		boids = append(boids, boid)
		points = append(points, *boid)
	}
	g.boids = boids
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
	r := float64(red) / 0xff
	g := float64(green) / 0xff
	b := float64(blue) / 0xff
	matrix.Translate(r, g, b, 0)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	b0, _ := g.getBoidById(0)
	b0pos := b0.position()
	re := rtreego.Point{b0pos.x, b0pos.y}.ToRect(separationRange)
	arr := rt.SearchIntersect(re)

	for _, boid := range g.boids {
		op.GeoM.Reset()
		op.ColorM.Reset()

		op.GeoM.Rotate(boid.angle)

		position := boid.position()
		op.GeoM.Translate(position.x, position.y)

		if boid.id == 0 {
			setRGB(&op.ColorM, 255, 0, 0)
		} else {
			for _, spa := range arr {
				b := spa.(Boid)
				if boid.id == b.id {
					setRGB(&op.ColorM, 0, 255, 0)
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
