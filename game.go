package main

import (
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Width   = 640
	Height  = 480
	fWidth  = float64(Width)
	fHeight = float64(Height)
)

var (
	boidCount   = 320
	trailLength = 40
)

type Game struct {
	boids  []*Boid
	tick   int
	pixels []byte
}

func createBoid(id int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()
	boidChan <- newBoid(id)
}

func updateBoid(boid *Boid, tick int, boidChan chan *Boid, wg *sync.WaitGroup) {
	defer wg.Done()
	boid.update(tick)
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	currentCount := len(g.boids)
	toCreateCount := boidCount - currentCount
	boidChan := make(chan *Boid, toCreateCount)
	wg.Add(boidCount)

	for i := 0; i < boidCount; i++ {
		if i >= currentCount {
			go createBoid(i, boidChan, &wg)
		} else {
			go updateBoid(g.boids[i], g.tick, boidChan, &wg)
		}
	}

	wg.Wait()
	close(boidChan)

	points := make([]rtreego.Spatial, 0, boidCount)
	for _, boid := range g.boids {
		points = append(points, *boid)
	}

	if toCreateCount > 0 {
		boids := make([]*Boid, 0, toCreateCount)
		for boid := range boidChan {
			boids = append(boids, boid)
			points = append(points, *boid)
		}
		g.boids = append(g.boids, boids...)
	}

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

	var wgt sync.WaitGroup
	wgt.Add(trailLength)
	for i := 0; i < trailLength; i++ {
		go func(trailPartIndex int) {
			defer wgt.Done()
			trailPosition := boid.trail[(g.tick+trailPartIndex)%trailLength]
			x := int(trailPosition.x)
			y := int(trailPosition.y)
			pixelDataPosition := (y*Width + x) * 4
			value := byte(255 * float64(trailLength-trailPartIndex) / float64(trailLength))
			g.pixels[pixelDataPosition] = value
			g.pixels[pixelDataPosition+1] = value
		}(i)
	}
	wgt.Wait()

	position := boid.Position()
	x := int(position.x)
	y := int(position.y)
	pixelDataPosition := (y*Width + x) * 4
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
	return Width, Height
}

func NewGame() *Game {
	boids := make([]*Boid, 0, boidCount)
	pixels := make([]byte, 4*Width*Height)
	return &Game{boids: boids, pixels: pixels}
}
