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
	width  = 320
	height = 240
)

var (
	img *ebiten.Image
)

func loadImage() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("boid.png")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	loadImage()

	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
