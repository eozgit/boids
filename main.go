package main

import (
	"context"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gonum.org/v1/gonum/num/quat"
)

var img *ebiten.Image
var rdb *redis.Client

var ctx = context.Background()

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("boid.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	boids []Boid
}

func (g *Game) Update() error {
	for _, boid := range g.boids {
		position := boid.getPosition()
		newPosition := position.Add(boid.velocity)
		boid.setPosition(newPosition)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	for _, boid := range g.boids {
		op := &ebiten.DrawImageOptions{}

		oa := boid.velocity.y / boid.velocity.x
		q := quat.Number{Real: oa}
		atan := quat.Atan(q)
		theta := atan.Real
		if boid.velocity.x > 0 {
			theta += math.Pi / 2
		} else {
			theta -= math.Pi / 2
		}
		op.GeoM.Rotate(theta)

		position := boid.getPosition()
		op.GeoM.Translate(position.x, position.y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	boids := []Boid{}
	for i := 0; i < 10; i++ {
		px := rand.Float64()*60 + 130
		py := rand.Float64()*60 + 90
		vx := rand.Float64() - .5
		vy := rand.Float64() - .5
		boid := Boid{
			id:       i,
			strId:    strconv.Itoa(i),
			velocity: &Vector{vx, vy},
		}
		boid.setPosition(&Vector{px, py})
		boids = append(boids, boid)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
