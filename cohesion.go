package main

type Cohesion struct{}

func (c *Cohesion) Delta(boid *Boid) (velocity *Vector) {
	velocity = &Vector{}
	neighbours, neighbourCount := GetNeighbours(boid.Position(), cohesionRange, boid.Id)
	if neighbourCount > 0 {
		neighbourPositions := &Vector{}
		for _, neighbour := range neighbours {
			neighbourPositions = neighbourPositions.Add(neighbour.Position())
		}
		velocity = neighbourPositions.Scale(1 / float64(neighbourCount)).Add(boid.Position().Negate()).Scale(cohesionWeight)
	}
	return
}
