package main

import (
	"fractals/pkg/gui"
	"fractals/pkg/screen"
	"log"
	"runtime"
)

const Width = 800
const Height = 800
const Title = "Hello world"

var home *screen.Home

func init() {
	runtime.LockOSThread()
}

func step(delta float64) {
	home.Update(delta)
	home.Render(delta)
}

func mousePress(x, y float64) {
	home.MousePress(x, y)
}

func main() {
	window, err := gui.NewWindow(Width, Height, Title)
	if err != nil {
		log.Fatal("Error opening window with error:", err)
	}

	home = new(screen.Home)
	home.Initialize(Width, Height)

	window.SetRunStep(step)
	window.SetMousePress(mousePress)
	window.Run()
}
