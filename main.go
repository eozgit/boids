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
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

var img *ebiten.Image
var rdb *redis.Client

var ctx = context.Background()

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("navigation.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Boid struct {
	id       int
	strId    string
	velocity *mat.Dense
}

func (boid *Boid) getPosition() *mat.Dense {
	pos := rdb.GeoPos(ctx, boid.strId, boid.strId).Val()[0]
	return mat.NewDense(1, 2, []float64{pos.Latitude * 10, pos.Longitude * 10})
}

type Game struct {
	boids []Boid
}

func (g *Game) Update() error {
	for _, boid := range g.boids {
		position := boid.getPosition()
		position.Add(position, boid.velocity)
		rdb.GeoAdd(ctx, boid.strId, &redis.GeoLocation{Name: boid.strId, Latitude: position.At(0, 0) / 10, Longitude: position.At(0, 1) / 10})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	for _, boid := range g.boids {
		op := &ebiten.DrawImageOptions{}

		oa := boid.velocity.At(0, 1) / boid.velocity.At(0, 0)
		q := quat.Number{Real: oa}
		atan := quat.Atan(q)
		theta := atan.Real
		if boid.velocity.At(0, 0) > 0 {
			theta += math.Pi / 2
		} else {
			theta -= math.Pi / 2
		}
		op.GeoM.Rotate(theta)

		op.GeoM.Scale(.02, .02)
		pos := boid.getPosition()
		px := pos.At(0, 0)
		py := pos.At(0, 1)
		op.GeoM.Translate(px, py)
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
			velocity: mat.NewDense(1, 2, []float64{vx, vy}),
		}
		cmd := rdb.GeoAdd(ctx, boid.strId, &redis.GeoLocation{Name: boid.strId, Latitude: px / 10, Longitude: py / 10})
		er := cmd.Err()
		if er != nil {
			log.Println(er)
		}
		boids = append(boids, boid)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
