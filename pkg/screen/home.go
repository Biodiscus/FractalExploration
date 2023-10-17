package screen

import (
	"fractals/pkg/gui"
	"fractals/pkg/mandelbrot"
	"image"
	"image/color"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Home struct {
	gui.Screen

	texture     uint32
	framebuffer uint32

	width  int
	height int

	xOffset  float64
	yOffset  float64
	origSize float64
	size     float64

	algo    *mandelbrot.MandelbrotGenerator
	context *image.RGBA

	recalcuate bool

	mousePressing bool
	mouseStartX   float64
	mouseStartY   float64
}

var Percision = uint64(5000)
var MovePercision = uint64(1000)

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

	h.algo = mandelbrot.NewGenerator(h.width, h.height, Percision, 2.0)
	// h.algo = mandelbrot.NewJuliaGenerator(h.width, h.height, Percision, 2.0)
	h.context = image.NewRGBA(image.Rect(0, 0, h.width, h.height))

	h.xOffset = -0.55
	h.yOffset = 0.0
	h.size = 2.0

	// h.xOffset = -0.19307385078394196
	// h.yOffset = 0.6484541110445733
	// h.size = 2.183869274407153e-08
	h.xOffset = -0.1930738568025074
	h.yOffset = 0.6484541050911828
	h.size = 4.046753855289677e-09

	// h.xOffset = -0.07153106897851419
	// h.yOffset = 0.8245921334837182
	// h.size = 0.00443706246892452
	h.origSize = 2.

	h.recalcuate = true
}

func (h *Home) MouseMove(x, y float64) {
	if !h.mousePressing {
		return
	}

	scaleFactor := h.origSize / h.size

	diffX := ((x - h.mouseStartX) / (250. * scaleFactor))
	diffY := ((y - h.mouseStartY) / (250. * scaleFactor))
	h.xOffset -= diffX
	h.yOffset += diffY

	h.recalcuate = true
	h.mouseStartX = x
	h.mouseStartY = y
}

func (h *Home) MousePress(action glfw.Action, x, y float64) {
	if action == glfw.Press {
		h.mousePressing = true
		h.mouseStartX = x
		h.mouseStartY = y

		h.algo.SetMaxIteration(MovePercision)
		h.recalcuate = true

	} else if action == glfw.Release {
		h.mousePressing = false

		diffX := x - h.mouseStartX
		diffY := y - h.mouseStartY
		// If there was no movement, zoom in instead
		if diffX == 0 && diffY == 0 {
			sizeDiff := h.size / h.origSize
			h.size -= 0.2 * sizeDiff
			Percision += 100
			MovePercision = Percision / 10

			log.Println("Offset X:", h.xOffset, ", offset Y:", h.yOffset, ", size:", h.size)
			log.Println("Percission:", Percision, "MovePercision", MovePercision)
		}

		h.algo.SetMaxIteration(uint64(Percision))
		h.recalcuate = true
	}
}

func (h *Home) Update(delta float64) {
	if h.recalcuate {
		amount := 100
		data := h.algo.GenerateImage(h.xOffset, h.yOffset, h.size, h.width, h.height, amount)

		for index, val := range data {
			x := index / h.width
			y := index % h.width
			col := color.RGBA{R: val.R, G: val.G, B: val.B, A: 255}
			h.context.Set(x, y, col)
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
