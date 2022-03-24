package main

import "github.com/dhconnelly/rtreego"

var rt *rtreego.Rtree

func createIndex(points ...rtreego.Spatial) {
	rt = rtreego.NewTree(2, 0, boidCount, points...)
}

func search(position *Vector, sideLength float64) []rtreego.Spatial {
	return rt.SearchIntersect(rtreego.Point{position.x, position.y}.ToRect(sideLength))
}
