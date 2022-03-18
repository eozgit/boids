package main

import (
	"image/color"
	"sync"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	homingWeight    = .01
	alignmentWeight = .002
	alignmentRadius = 30.
	op              = &ebiten.DrawImageOptions{}
)

type Game struct {
	boids []*Boid
}

func (g *Game) Update() error {
	var wg sync.WaitGroup
	pointChan := make(chan kdbush.Point, boidCount)
	for _, boid := range g.boids {
		wg.Add(1)
		var wgv sync.WaitGroup
		velocityChan := make(chan *Vector, 2)
		go func(b *Boid) {
			defer wg.Done()
			wgv.Add(2)
			go getHomingVelocity(b.position, velocityChan, &wgv)
			go getAlignmentVelocity(b, velocityChan, &wgv)
			wgv.Wait()
			close(velocityChan)
			for velocity := range velocityChan {
				b.velocity.Add(velocity)
			}
			b.velocity.Limit(.5)
			b.position.Add(b.velocity)
			b.calculateAngle()
			pointChan <- &kdbush.SimplePoint{X: b.position.x, Y: b.position.y}
		}(boid)
	}
	wg.Wait()
	close(pointChan)
	points = []kdbush.Point{}
	for point := range pointChan {
		points = append(points, point)
	}
	bush = kdbush.NewBush(points, boidCount)
	return nil
}

func getAlignmentVelocity(boid *Boid, chanVelocity chan *Vector, wgv *sync.WaitGroup) {
	defer wgv.Done()
	arr := bush.Within(&kdbush.SimplePoint{X: boid.position.x, Y: boid.position.y}, alignmentRadius)
	vector := &Vector{}
	for _, i := range arr {
		neighbour := boids[i]
		if neighbour.id != boid.id {
			vector.Add(neighbour.velocity)
		}
	}
	l := float64(len(arr))
	vector.x = vector.x / l
	vector.y = vector.y / l
	vector.Scale(alignmentWeight)
	chanVelocity <- vector
}

func getHomingVelocity(position *Vector, chanVelocity chan *Vector, wgv *sync.WaitGroup) {
	defer wgv.Done()
	x, y := 0., 0.
	fWidth, fHeight := float64(width), float64(height)
	outOfBoundsLeft := position.x < 0
	outOfBoundsRight := position.x > fWidth
	outOfBoundsTop := position.y < 0
	outOfBoundsBottom := position.y > fHeight
	if outOfBoundsLeft {
		x = -position.x / fWidth
	} else if outOfBoundsRight {
		x = -(position.x - fWidth) / fWidth
	}
	if outOfBoundsTop {
		y = -position.y / fHeight
	} else if outOfBoundsBottom {
		y = -(position.y - fHeight) / fHeight
	}
	vector := &Vector{x, y}
	vector.Scale(homingWeight)
	chanVelocity <- vector
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	for _, boid := range g.boids {
		op.GeoM.Reset()

		op.GeoM.Rotate(boid.angle)

		op.GeoM.Translate(boid.position.x, boid.position.y)
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
