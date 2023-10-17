package mandelbrot

import (
	"log"
	"sync"
	"time"

	"github.com/essentialkaos/ek/v12/color"
)

func (m *MandelbrotGenerator) GenerateImage(xOffset, yOffset, size float64, width, height int, amount int) []color.RGB {
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
			go m.calclulateBlock(xOffset, yOffset, size, width, rect, data, &wg)
		}
	}

	wg.Wait()

	end := time.Now().UnixMilli()
	log.Println("[Percission:", m.max_iteration, "]Generation took: ", end-start, "ms")

	return data
}

func (m *MandelbrotGenerator) calclulateBlock(xOffset, yOffset, size float64, originalWidth int, bounds Rectangle, data []color.RGB, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := bounds.Left; x < bounds.Right; x++ {
		for y := bounds.Bottom; y < bounds.Top; y++ {
			data[originalWidth*x+y] = m.Color(xOffset, yOffset, size, float64(x), float64(y))
		}
	}
}
