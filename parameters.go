package main

const (
	maxVel      = .5
	trailLength = 40
)

type Parameters struct {
	maxVel      *FloatParam
	trailLength *IntParam
}

func newParameters() *Parameters {
	return &Parameters{newFloatParam(maxVel), newIntParam(trailLength)}
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
	p.field--
}

func newIntParam(value int) *IntParam {
	return &IntParam{value}
}

var _ Param[int] = (*IntParam)(nil)

type FloatParam struct {
	field float64
}

func (p *FloatParam) value() float64 {
	return p.field
}

func (p *FloatParam) increase() {
	p.field++
}

func (p *FloatParam) decrease() {
	p.field--
}

func newFloatParam(value float64) *FloatParam {
	return &FloatParam{value}
}

var _ Param[float64] = (*FloatParam)(nil)
