package main

import (
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width     = 320
	height    = 240
	boidCount = 25
)

var (
	img    *ebiten.Image
	boids  []*Boid
	points []kdbush.Point
	bush   *kdbush.KDBush
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
	velocity := rand.Float64() - .5
	if math.Abs(velocity) < .2 {
		velocity *= 2
	}
	return velocity
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
				strId:    strconv.Itoa(id),
				position: &Vector{px, py},
				velocity: &Vector{vx, vy},
			}
			boid.calculateAngle()
			boidChan <- boid
		}(i)
	}
	wg.Wait()
	close(boidChan)
	boids = []*Boid{}
	points = []kdbush.Point{}
	for boid := range boidChan {
		boids = append(boids, boid)
		points = append(points, &kdbush.SimplePoint{X: boid.position.x, Y: boid.position.y})
	}
	bush = kdbush.NewBush(points, boidCount)

	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
