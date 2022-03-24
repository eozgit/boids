package main

import (
	"errors"
	"image/color"
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	boidCount = 100
	op        = &ebiten.DrawImageOptions{}
)

type Game struct {
	boids []*Boid
}

func createBoid(id int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()
	px := rand.Float64() * width
	py := rand.Float64() * height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5
	boid := &Boid{
		id:       id,
		Point:    rtreego.Point{px, py},
		velocity: &Vector{vx, vy},
	}
	boid.calculateAngle()
	boidChan <- boid
}

func updateBoid(b *Boid, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()

	b.calculateVelocity()

	position := b.position().Add(b.velocity)
	Wrap(position)
	b.Point = rtreego.Point{position.x, position.y}
	b.calculateAngle()
	boidChan <- b
}

func Wrap(position *Vector) {
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

func (g *Game) Update() error {
	var wg sync.WaitGroup
	boidChan := make(chan *Boid, boidCount)
	wg.Add(boidCount)

	actualCount := len(g.boids)
	for i := 0; i < boidCount; i++ {
		if i >= actualCount {
			go createBoid(i, boidChan, &wg)
		} else {
			boid := g.boids[i]
			go updateBoid(boid, boidChan, &wg)
		}
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
	createIndex(points...)

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
	arr := search(b0.position(), separationRange)

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
