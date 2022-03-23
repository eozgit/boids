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
	boidCount = 10
	op        = &ebiten.DrawImageOptions{}
	rt        *rtreego.Rtree
)

type Game struct {
	boids []*Boid
	rt    *rtreego.Rtree
}

func createBoid(id int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()
	px := randPostition(width)
	py := randPostition(height)
	vx := randVelocity()
	vy := randVelocity()
	boid := &Boid{
		id:       id,
		velocity: &Vector{vx, vy, "vel"},
		Point:    rtreego.Point{px, py},
	}
	boid.calculateAngle()
	boidChan <- boid
}

func updateBoid(b *Boid, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()

	velocityCalc := VelocityCalculator{}
	velocityCalc.calculate(b)

	position := b.position()
	position.Add(b.velocity)
	b.Point = rtreego.Point{position.x, position.y}
	b.calculateAngle()
	boidChan <- b
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

func randPostition(dim float64) float64 {
	position := (rand.Float64() - .5) * dim
	if position > 0 {
		position += dim
	}
	return position
}

func randVelocity() float64 {
	return rand.Float64() - .5
}
