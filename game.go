package main

import (
	"math/rand"
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	boidCount   = 120
	trailLength = 50
)

type Game struct {
	boids  []*Boid
	tick   int
	pixels []byte
}

func createBoid(id int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()
	px := rand.Float64() * width
	py := rand.Float64() * height
	vx := rand.Float64() - .5
	vy := rand.Float64() - .5
	trail := make([]Vector, trailLength)
	for i := 0; i < trailLength; i++ {
		trail = append(trail, Vector{px, py})
	}
	boid := &Boid{
		id:       id,
		Point:    rtreego.Point{px, py},
		velocity: &Vector{vx, vy},
		trail:    trail,
	}
	boidChan <- boid
}

func updateBoid(boid *Boid, tick int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()

	position := boid.position()
	boid.trail[tick%trailLength] = *position

	boid.calculateVelocity()

	position = position.Add(boid.velocity)
	wrap(position)
	boid.Point = rtreego.Point{position.x, position.y}
	boidChan <- boid
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
			go updateBoid(boid, g.tick, boidChan, &wg)
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

	g.tick++

	return nil
}

func (g *Game) resetPixels() {
	for i := range g.pixels {
		g.pixels[i] = 255
	}
}

func (g *Game) drawBoid(boid *Boid, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < trailLength; i++ {
		trailPosition := boid.trail[(g.tick+i)%trailLength]
		x := int(trailPosition.x)
		y := int(trailPosition.y)
		pixelDataPosition := (y*width + x) * 4
		value := byte(255 * float64(trailLength-i) / float64(trailLength))
		g.pixels[pixelDataPosition] = value
		g.pixels[pixelDataPosition+1] = value
	}
	position := boid.position()
	x := int(position.x)
	y := int(position.y)
	pixelDataPosition := (y*width + x) * 4
	g.pixels[pixelDataPosition] = 0
	g.pixels[pixelDataPosition+1] = 0
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.resetPixels()

	var wg sync.WaitGroup
	wg.Add(len(g.boids))
	for _, boid := range g.boids {
		go g.drawBoid(boid, &wg)
	}
	wg.Wait()

	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func NewGame() *Game {
	pixels := make([]byte, 4*width*height)
	return &Game{pixels: pixels}
}
