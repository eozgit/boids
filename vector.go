package main

import "gonum.org/v1/gonum/mat"

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Add(v2 *Vector) *Vector {
	matrix := mat.NewDense(1, 2, []float64{v.x, v.y})
	matrix2 := mat.NewDense(1, 2, []float64{v2.x, v2.y})
	matrix.Add(matrix, matrix2)
	return &Vector{matrix.At(0, 0), matrix.At(0, 1)}
}
