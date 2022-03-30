package main

import (
	"github.com/gravestench/mathlib"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) checkInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		position := mathlib.Vector2{X: float64(x), Y: float64(y)}
		g.addBoid(&position)
	}
}
