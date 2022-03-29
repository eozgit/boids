package main

import (
	"github.com/dhconnelly/rtreego"
	"github.com/gravestench/mathlib"
)

type Spatial struct {
	tree *rtreego.Rtree
}

func (i *Spatial) search(position *mathlib.Vector2, sideLength float64) []rtreego.Spatial {
	return i.tree.SearchIntersect(rtreego.Point{position.X, position.Y}.ToRect(sideLength))
}

func (i *Spatial) GetNeighbours(position *mathlib.Vector2, sideLength float64, boidId int) ([]Boid, int) {
	var neighbours []Boid
	spatials := i.search(position, sideLength)
	for _, spatial := range spatials {
		neighbour := spatial.(Boid)
		if boidId == neighbour.Id {
			continue
		}

		neighbours = append(neighbours, neighbour)
	}

	return neighbours, len(neighbours)
}

func newIndex(points ...rtreego.Spatial) *Spatial {
	return &Spatial{
		tree: rtreego.NewTree(2, 0, len(points), points...),
	}
}
