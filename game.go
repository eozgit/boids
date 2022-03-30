package main

import (
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/gravestench/mathlib"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Width   = 640
	Height  = 480
	fWidth  = float64(Width)
	fHeight = float64(Height)
)

type Game struct {
	boidCount int
	boids     []*Boid
	tick      int
	pixels    []byte
}

func (g *Game) Update() error {
	g.checkInput()
	var wg sync.WaitGroup
	currentCount := len(g.boids)
	toCreateCount := g.boidCount - currentCount
	boidChan := make(chan *Boid, toCreateCount)
	wg.Add(g.boidCount)

	for i := 0; i < g.boidCount; i++ {
		idx := i
		go func() {
			defer wg.Done()
			if idx >= currentCount {
				boidChan <- newBoid(idx, nil)
			} else {
				boid := g.boids[idx]
				boid.update(g.tick)
			}
		}()
	}

	wg.Wait()
	close(boidChan)

	points := make([]rtreego.Spatial, 0, g.boidCount)
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

	global.setIndex(newIndex(points...))

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

	trailChan := make(chan *TrailPixel, global.params.trailLength.value())

	boid.getTrailPixels(g.tick, trailChan)

	for trailPixel := range trailChan {
		g.pixels[trailPixel.pixelIndex] = trailPixel.colourValue
		g.pixels[trailPixel.pixelIndex+1] = trailPixel.colourValue
	}

	position := boid.Position()
	x := int(position.X)
	y := int(position.Y)
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

func (g *Game) addBoid(position *mathlib.Vector2) {
	boid := newBoid(g.boidCount, position)
	g.boids = append(g.boids, boid)
	g.boidCount++
}

func NewGame() *Game {
	boidCount := 200
	boids := make([]*Boid, 0, boidCount)
	pixels := make([]byte, 4*Width*Height)
	return &Game{boidCount: boidCount, boids: boids, pixels: pixels}
}
