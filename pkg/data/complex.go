package data

import "math"

type Complex struct {
	Real      float64
	Imaginary float64
}

func NewComplex(Real, Imaginary float64) *Complex {
	complex := Complex{
		Real:      Real,
		Imaginary: Imaginary,
	}
	return &complex
}

func (m *Complex) Abs() float64 {
	return math.Hypot(m.Real, m.Imaginary)
}

func (m *Complex) Phase() float64 {
	return math.Atan2(m.Imaginary, m.Real)
}

func (m *Complex) Plus(c Complex) *Complex {
	Real := m.Real + c.Real
	imag := m.Imaginary + c.Imaginary
	return NewComplex(Real, imag)
}

func (m *Complex) Minus(c Complex) *Complex {
	Real := m.Real - c.Real
	imag := m.Imaginary - c.Imaginary
	return NewComplex(Real, imag)
}

func (m *Complex) Times(c Complex) *Complex {
	Real := m.Real*c.Real - m.Imaginary*c.Imaginary
	imag := m.Real*c.Imaginary + m.Imaginary*c.Real
	return NewComplex(Real, imag)
}

func (m *Complex) Scale(alpha float64) *Complex {
	return NewComplex(
		alpha*m.Real,
		alpha*m.Imaginary,
	)
}
