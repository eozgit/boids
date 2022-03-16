package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gonum.org/v1/gonum/num/quat"
)

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

type Game struct {
	boids []Boid
}

func (g *Game) Update() error {
	points = []kdbush.Point{}
	for _, boid := range g.boids {
		newPosition := boid.position.Add(boid.velocity)
		boid.setPosition(newPosition)
	}
	bush = kdbush.NewBush(points, len(g.boids))
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

		op.GeoM.Translate(boid.position.x, boid.position.y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	rand.Seed(time.Now().UnixNano())

	points = []kdbush.Point{}

	boids := []Boid{}
	for i := 0; i < 10; i++ {
		px := rand.Float64()*60 + 130
		py := rand.Float64()*60 + 90
		vx := rand.Float64() - .5
		vy := rand.Float64() - .5
		boid := Boid{
			id:       i,
			strId:    strconv.Itoa(i),
			position: &Vector{},
			velocity: &Vector{vx, vy},
		}
		boid.setPosition(&Vector{px, py})
		boids = append(boids, boid)
	}
	bush = kdbush.NewBush(points, len(boids))

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
