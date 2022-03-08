package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("navigation.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Boid struct {
	position *mat.Dense
	velocity *mat.Dense
}

type Game struct {
	boids []Boid
}

func (g *Game) Update() error {
	for _, boid := range g.boids {
		boid.position.Add(boid.position, boid.velocity)
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
		op.GeoM.Translate(boid.position.At(0, 0), boid.position.At(0, 1))
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	rand.Seed(time.Now().UnixNano())

	boids := []Boid{}
	for i := 0; i < 10; i++ {
		boid := Boid{
			position: mat.NewDense(1, 2, []float64{rand.Float64()*60 + 130, rand.Float64()*60 + 90}),
			velocity: mat.NewDense(1, 2, []float64{rand.Float64() - .5, rand.Float64() - .5}),
		}
		boids = append(boids, boid)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{boids: boids}); err != nil {
		log.Fatal(err)
	}
}
