package main

type Global struct {
	params                 *Parameters
	index                  *Spatial
	velocityComponents     []Velocity
	velocityComponentCount int
}

func (g *Global) setIndex(index *Spatial) {
	g.index = index
}

func newGlobal() *Global {
	velocityComponents := []Velocity{
		newSeparation(),
		newAlignment(),
		newCohesion(),
		newNoise(),
	}

	return &Global{
		params:                 newParameters(),
		index:                  newIndex(),
		velocityComponents:     velocityComponents,
		velocityComponentCount: len(velocityComponents),
	}
}
