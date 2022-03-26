package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(Width*2, Height*2)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
