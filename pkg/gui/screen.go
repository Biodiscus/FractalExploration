package gui

import "github.com/go-gl/glfw/v3.2/glfw"

type Screen interface {
	Initialize()
	MousePress(action glfw.Action, x, y float64)
	MouseMove(x, y float64)
	Update(delta float64)
	Render(delta float64)
}
