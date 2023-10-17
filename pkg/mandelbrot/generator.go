package mandelbrot

import (
	"fractals/pkg/data"
	"math"
	"math/cmplx"

	"github.com/essentialkaos/ek/v12/color"
)

type MandelbrotGenerator struct {
	escapeRadius  float64
	max_iteration uint64

	width  float64
	height float64
}

func NewGenerator(width, height int, max_iteration uint64, escapeRadius float64) *MandelbrotGenerator {
	m := MandelbrotGenerator{}
	m.escapeRadius = escapeRadius
	m.width = float64(width)
	m.height = float64(height)
	m.max_iteration = max_iteration
	return &m
}

func (m *MandelbrotGenerator) Generate(xOffset, yOffset, size float64) [][]color.RGB {
	width := int(m.width)
	height := int(m.height)

	values := [][]color.RGB{}
	for x := 0; x < width; x++ {

		values = append(values, []color.RGB{})

		for y := 0; y < height; y++ {
			data := m.Color(
				float64(xOffset), float64(yOffset), float64(size),
				float64(x), float64(y),
			)
			values[x] = append(values[x], data)
		}
	}

	return values
}

func (m *MandelbrotGenerator) Color(x float64, y float64, size float64, x_iteration float64, y_iteration float64) color.RGB {
	data := m.data(x, y, size, x_iteration, y_iteration)
	calc := data * 2.0 / float64(m.max_iteration)
	// log.Println("Data: ", data, "calc:", calc)
	brightness := 1.
	if floatEquals(calc, 2.) {
		brightness = 0.
	}

	return color.HSV2RGB(color.HSV{H: calc, S: 1.0, V: brightness})
}

func (m *MandelbrotGenerator) data(x float64, y float64, size float64, x_iteration float64, y_iteration float64) float64 {
	x0 := x - size/2 + size*x_iteration/m.width
	y0 := y - size/2 + size*y_iteration/m.height

	z := data.NewComplex(x0, y0)

	return m.step(*z, m.max_iteration)
}

func (m *MandelbrotGenerator) step(z0 data.Complex, max uint64) float64 {
	zz0 := complex(z0.Real, z0.Imaginary)
	zz := zz0

	for i := uint64(0); i < max; i++ {
		abs := cmplx.Abs(zz)

		if abs > 2 {
			sub := (math.Log(abs) / math.Log(m.escapeRadius))
			return float64(i) + 1.0 - float64(sub)
		}

		zz = (zz * zz) + zz0
	}

	return float64(max)
}

func floatEquals(f1, f2 float64) bool {
	return floatEqualsPercission(f1, f2, 0.00000001)
}

func floatEqualsPercission(f1, f2, percission float64) bool {
	return math.Abs(f1-f2) < float64(percission)
}

func (m *MandelbrotGenerator) SetMaxIteration(iteration uint64) {
	m.max_iteration = iteration
}

func (m *MandelbrotGenerator) SetEscapeRadius(radius float64) {
	m.escapeRadius = radius
}
