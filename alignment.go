package main

type Alignment struct{}

var _ Velocity = (*Alignment)(nil)

func (_ *Alignment) Delta(boid *Boid) *Vector {
	velocity := &Vector{}
	neighbours, neighbourCount := GetNeighbours(boid.Position(), boid.params.alignmentRange.value(), boid.Id)
	if neighbourCount == 0 {
		return velocity
	}

	for _, neighbour := range neighbours {
		velocity = velocity.Add(neighbour.Velocity)
	}

	return velocity.Scale(1 / float64(neighbourCount)).Scale(boid.params.alignmentWeight.value())
}
