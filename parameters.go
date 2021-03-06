package main

const (
	step             = .01
	maximumVelocity  = .5
	separationRange  = 6.
	separationWeight = .02
	alignmentRange   = 19.
	alignmentWeight  = .01
	cohesionRange    = 19.
	cohesionWeight   = .0004
	noiseWeight      = .03
)

type Parameters struct {
	maximumVelocity  *FloatParam
	separationRange  *FloatParam
	separationWeight *FloatParam
	alignmentRange   *FloatParam
	alignmentWeight  *FloatParam
	cohesionRange    *FloatParam
	cohesionWeight   *FloatParam
	noiseWeight      *FloatParam
	trailLength      *IntParam
}

func newParameters() *Parameters {
	return &Parameters{
		maximumVelocity:  newFloatParam(maximumVelocity),
		separationRange:  newFloatParam(separationRange),
		separationWeight: newFloatParam(separationWeight),
		alignmentRange:   newFloatParam(alignmentRange),
		alignmentWeight:  newFloatParam(alignmentWeight),
		cohesionRange:    newFloatParam(cohesionRange),
		cohesionWeight:   newFloatParam(cohesionWeight),
		noiseWeight:      newFloatParam(noiseWeight),
		trailLength:      newIntParam(trailLength),
	}
}

type Param[T int | float64] interface {
	value() T
	increase()
	decrease()
}

type IntParam struct {
	field int
}

func (p *IntParam) value() int {
	return p.field
}

func (p *IntParam) increase() {
	p.field++
}

func (p *IntParam) decrease() {
	if p.field > 0 {
		p.field--
	}
}

func newIntParam(value int) *IntParam {
	return &IntParam{value}
}

var _ Param[int] = (*IntParam)(nil)

type FloatParam struct {
	base   float64
	factor float64
}

func (p *FloatParam) value() float64 {
	return p.base * p.factor
}

func (p *FloatParam) increase() {
	p.factor += step
}

func (p *FloatParam) decrease() {
	if p.factor > step {
		p.factor -= step
	}
}

func newFloatParam(value float64) *FloatParam {
	return &FloatParam{value, 1}
}

var _ Param[float64] = (*FloatParam)(nil)
