package screen

import (
	"fractals/pkg/gui"
	"fractals/pkg/mandelbrot"
	"image"
	"image/color"

	"github.com/go-gl/gl/v2.1/gl"
)

type Home struct {
	gui.Screen

	texture     uint32
	framebuffer uint32

	width  int
	height int

	xOffset float64
	yOffset float64
	size    float64

	algo    *mandelbrot.MandelbrotGenerator
	context *image.RGBA

	recalcuate bool
}

func (h *Home) Initialize(width, height int) {
	h.width = width
	h.height = height

	{ // Texture
		gl.GenTextures(1, &h.texture)
		gl.BindTexture(gl.TEXTURE_2D, h.texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		// gl.BindImageTexture(0, h.texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
	}
	{ // Framebuffer
		gl.GenFramebuffers(1, &h.framebuffer)
		gl.BindFramebuffer(gl.FRAMEBUFFER, h.framebuffer)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, h.texture, 0)

		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, h.framebuffer)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	}

	h.algo = mandelbrot.NewGenerator(h.width, h.height, 100, 2.)
	h.context = image.NewRGBA(image.Rect(0, 0, h.width, h.height))

	h.xOffset = -0.5
	h.yOffset = 0.0
	h.size = 2.0

	h.recalcuate = true
}

func (h *Home) MousePress(x, y float64) {
	h.size -= 0.1

	// coord = mv + ( 4*value/win_size - 2 )/zoom_factor
	// cooord.x = mv.x + (4.0 value / win_size.width - 2) / zoom
	// cooord.y = mv.y + (4.0 value / win_size.height - 2) / zoom

	// coordX := x + (4.0*h.xOffset/float64(h.width)-2.)/h.size
	// coordY := y + (4.0*h.yOffset/float64(h.height)-2.)/h.size
	// log.Println("X:", coordX, "Y:", coordY)

	// coordsX := x - (float64(h.width) / 2.0)
	// coordsY := y - (float64(h.height) / 2.0)

	// log.Println("CoordsX:", coordsX, "CoordsY:", coordsY)

	h.recalcuate = true
}

func (h *Home) Update(delta float64) {
	if h.recalcuate {
		for x := 0; x < h.width; x++ {
			for y := 0; y < h.height; y++ {
				col := h.algo.Color(h.xOffset, h.yOffset, h.size, float64(x), float64(y))
				h.context.Set(x, y, color.RGBA{R: col.R, G: col.G, B: col.B, A: 255})
			}
		}
		h.recalcuate = false
	}
}

func (h *Home) Render(delta float64) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, h.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(h.width), int32(h.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(h.context.Pix))

	width := int32(h.width)
	height := int32(h.height)

	gl.BlitFramebuffer(0, 0, width, height, 0, 0, width*2, height*2, gl.COLOR_BUFFER_BIT, gl.LINEAR)
}
