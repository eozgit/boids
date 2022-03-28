package main

import (
	"github.com/dhconnelly/rtreego"
)

var tree *rtreego.Rtree

func createIndex(points ...rtreego.Spatial) {
	tree = rtreego.NewTree(2, 0, len(points), points...)
}

func search(position *Vector, sideLength float64) []rtreego.Spatial {
	return tree.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(sideLength))
}

func GetNeighbours(position *Vector, sideLength float64, boidId int) ([]Boid, int) {
	var neighbours []Boid
	spatials := search(position, sideLength)
	for _, spatial := range spatials {
		neighbour := spatial.(Boid)
		if boidId == neighbour.Id {
			continue
		}

		neighbours = append(neighbours, neighbour)
	}

	return neighbours, len(neighbours)
}
