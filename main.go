package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var boidCount = 10

var img *ebiten.Image
var points []kdbush.Point
var bush *kdbush.KDBush

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("boid.png")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	boidChan := make(chan Boid, boidCount)

	var wg sync.WaitGroup
	for i := 0; i < boidCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			px := rand.Float64()*60 + 130
			py := rand.Float64()*60 + 90
			vx := rand.Float64() - .5
			vy := rand.Float64() - .5
			boid := Boid{
				id:       id,
				strId:    strconv.Itoa(id),
				position: &Vector{},
				velocity: &Vector{vx, vy},
			}
			boid.setPosition(&Vector{px, py})
			boidChan <- boid
		}(i)
	}
	wg.Wait()
	close(boidChan)
	boids := []Boid{}
	points = []kdbush.Point{}
	for boid := range boidChan {
		boids = append(boids, boid)
		points = append(points, &kdbush.SimplePoint{X: boid.position.x, Y: boid.position.y})
	}
	bush = kdbush.NewBush(points, boidCount)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
