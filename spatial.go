package main

import (
	"github.com/dhconnelly/rtreego"
)

var tree *rtreego.Rtree

func createIndex(points ...rtreego.Spatial) {
	tree = rtreego.NewTree(2, 0, boidCount, points...)
}

func search(position *Vector, sideLength float64) []rtreego.Spatial {
	return tree.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(sideLength))
}

func GetNeighbours(position *Vector, sideLength float64, boidId int) (neighbours []Boid, neighbourCount int) {
	spatials := search(position, sideLength)
	for _, spatial := range spatials {
		potential := spatial.(Boid)
		if boidId != potential.Id {
			neighbours = append(neighbours, potential)
		}
	}
	neighbourCount = len(neighbours)
	return
}
