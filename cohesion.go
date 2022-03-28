package main

type Cohesion struct{}

var _ Velocity = (*Cohesion)(nil)

func (_ *Cohesion) Delta(boid *Boid) *Vector {
	velocity := &Vector{}
	neighbours, neighbourCount := GetNeighbours(boid.Position(), boid.params.cohesionRange.value(), boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	neighbourPositions := &Vector{}
	for _, neighbour := range neighbours {
		neighbourPositions = neighbourPositions.Add(neighbour.Position())
	}

	return neighbourPositions.Scale(1 / float64(neighbourCount)).Add(boid.Position().Negate()).Scale(boid.params.cohesionWeight.value())
}
