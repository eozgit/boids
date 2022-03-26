package main

type Alignment struct{}

func (a *Alignment) Delta(boid *Boid) (velocity *Vector) {
	velocity = &Vector{}
	neighbours, neighbourCount := GetNeighbours(boid.Position(), alignmentRange, boid.Id)
	if neighbourCount > 0 {
		for _, neighbour := range neighbours {
			velocity = velocity.Add(neighbour.Velocity)
		}
		velocity = velocity.Scale(1 / float64(neighbourCount)).Scale(alignmentWeight)
	}
	return
}
