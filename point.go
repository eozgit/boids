package main

import "github.com/dhconnelly/rtreego"

var tol = 0.01

type Point struct {
	loc  rtreego.Point
	boid *Boid
}

func (s Point) Bounds() *rtreego.Rect {
	return s.loc.ToRect(tol)
}
