package main

type Cohesion struct{}

var _ Velocity = (*Cohesion)(nil)

func (_ *Cohesion) Delta(boid *Boid) *Vector {
	velocity := &Vector{}
	neighbours, neighbourCount := global.index.GetNeighbours(boid.Position(), global.params.cohesionRange.value(), boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	neighbourPositions := &Vector{}
	for _, neighbour := range neighbours {
		neighbourPositions = neighbourPositions.Add(neighbour.Position())
	}

	return neighbourPositions.Scale(1 / float64(neighbourCount)).Add(boid.Position().Negate()).Scale(global.params.cohesionWeight.value())
}

func newCohesion() *Cohesion {
	return &Cohesion{}
}
