package gui

type Screen interface {
	Initialize()
	MousePress(x, y float64)
	Update(delta float64)
	Render(delta float64)
}
