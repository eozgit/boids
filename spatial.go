package main

import (
	"github.com/dhconnelly/rtreego"
)

type Spatial struct {
	tree *rtreego.Rtree
}

func (i *Spatial) search(position *Vector, sideLength float64) []rtreego.Spatial {
	return i.tree.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(sideLength))
}

func (i *Spatial) GetNeighbours(position *Vector, sideLength float64, boidId int) ([]Boid, int) {
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
