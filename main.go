package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width   = 640
	height  = 480
	fWidth  = float64(width)
	fHeight = float64(height)
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
