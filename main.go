package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width     = 320
	height    = 240
	boidCount = 50
)

var (
	img *ebiten.Image
	rt  *rtreego.Rtree
)

func loadImage() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("boid.png")
	if err != nil {
		log.Fatal(err)
	}
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

func main() {
	rand.Seed(time.Now().UnixNano())

	loadImage()

	boidChan := make(chan *Boid, boidCount)

	var wg sync.WaitGroup
	for i := 0; i < boidCount; i++ {
		wg.Add(1)
		go func(id int) {
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
		}(i)
	}
	wg.Wait()
	close(boidChan)
	boids := []*Boid{}
	points := []rtreego.Spatial{}
	for boid := range boidChan {
		boids = append(boids, boid)
		points = append(points, *boid)
	}
	rt = rtreego.NewTree(2, 5, 500, points...)

	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
