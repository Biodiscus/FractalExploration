package main

import (
	"fmt"
	"fractals/pkg/mandelbrot"
	"image"
	"image/color"
	"image/png"
	"os"

	ek12color "github.com/essentialkaos/ek/v12/color"
)

func main() {
	w := 500
	h := 500

	gen := mandelbrot.NewGenerator(w, h, 200, 3.)

	for size := 2.0; size > 0.1; size -= 0.1 {
		data := gen.Generate(-0.5, 0, size)
		createImage(w, h, data, fmt.Sprintf("size-%f.png", size))
	}
}

func createImage(width, height int, data [][]ek12color.RGB, title string) {
	top_left := image.Point{0, 0}
	bot_right := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{top_left, bot_right})

	for x := 0; x < len(data); x++ {
		for y := 0; y < len(data[x]); y++ {
			col := data[x][y]
			img.Set(x, y, color.RGBA{R: col.R, G: col.G, B: col.B, A: 255})
		}
	}

	imgFile, _ := os.Create("images/" + title)
	png.Encode(imgFile, img)
}
