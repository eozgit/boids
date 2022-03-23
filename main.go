package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width   = 320
	height  = 240
	fWidth  = float64(width)
	fHeight = float64(height)
)

var (
	img *ebiten.Image
)

func loadImage() (err error) {
	img, _, err = ebitenutil.NewImageFromFile("boid.png")
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := loadImage(); err != nil {
		log.Fatal(err)
		return
	}

	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
