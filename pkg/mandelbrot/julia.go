package mandelbrot

import (
	"fractals/pkg/data"
	"log"
	"math"
	"math/cmplx"
	"sync"
	"time"

	"github.com/essentialkaos/ek/v12/color"
)

type JuliaGenerator struct {
	escapeRadius  float64
	max_iteration uint64

	width  float64
	height float64
}

func NewJuliaGenerator(width, height int, max_iteration uint64, escapeRadius float64) *JuliaGenerator {
	j := JuliaGenerator{}
	j.escapeRadius = escapeRadius
	j.width = float64(width)
	j.height = float64(height)
	j.max_iteration = max_iteration
	return &j
}

func (j *JuliaGenerator) GenerateImage(xOffset, yOffset, size float64, width, height int, amount int) []color.RGB {
	var wg sync.WaitGroup

	newWidth := width / amount
	newHeight := height / amount
	data := make([]color.RGB, width*height)

	start := time.Now().UnixMilli()

	for x := 0; x < newWidth; x++ {
		xPos := x * amount

		for y := 0; y < newHeight; y++ {
			yPos := y * amount
			rect := Rectangle{xPos, yPos, xPos + amount, yPos + amount}

			wg.Add(1)
			go j.calclulateBlock(xOffset, yOffset, size, width, rect, data, &wg)
		}
	}

	wg.Wait()

	end := time.Now().UnixMilli()
	log.Println("[Percission:", j.max_iteration, "]Generation took: ", end-start, "ms")

	return data
}

func (j *JuliaGenerator) calclulateBlock(xOffset, yOffset, size float64, originalWidth int, bounds Rectangle, data []color.RGB, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := bounds.Left; x < bounds.Right; x++ {
		for y := bounds.Bottom; y < bounds.Top; y++ {
			data[originalWidth*x+y] = j.Color(xOffset, yOffset, size, float64(x), float64(y))
		}
	}
}

func (j *JuliaGenerator) Color(x float64, y float64, size float64, x_iteration float64, y_iteration float64) color.RGB {
	data := j.data(x, y, size, x_iteration, y_iteration)
	calc := data * 2.0 / float64(j.max_iteration)
	brightness := 1.
	if floatEquals(calc, 2.) {
		brightness = 0.
	}

	return color.HSV2RGB(color.HSV{H: calc, S: 1.0, V: brightness})
}

func (j *JuliaGenerator) data(x float64, y float64, size float64, x_iteration float64, y_iteration float64) float64 {
	x0 := x - size/2 + size*x_iteration/j.width
	y0 := y - size/2 + size*y_iteration/j.height

	z := data.NewComplex(x0, y0)
	c := data.NewComplex(-0.75, 0.1)

	return j.step(*c, *z, j.max_iteration)
}

func (j *JuliaGenerator) step(c data.Complex, z0 data.Complex, max uint64) float64 {
	zz := complex(z0.Real, z0.Imaginary)
	cc := complex(c.Real, c.Imaginary)

	for i := uint64(0); i < max; i++ {
		abs := cmplx.Abs(zz)

		if abs > 2 {
			sub := (math.Log(abs) / math.Log(j.escapeRadius))
			return float64(i) + 1.0 - float64(sub)
		}

		zz = (zz * zz) + cc
	}

	return float64(max - 1)
}

func (j *JuliaGenerator) SetMaxIteration(iteration uint64) {
	j.max_iteration = iteration
}

func (j *JuliaGenerator) SetEscapeRadius(radius float64) {
	j.escapeRadius = radius
}
