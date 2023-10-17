package gui

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	width  int
	height int

	title string

	previousTime float64

	step       RunCallback
	mousePress MousePressCallback
	mouseMove  MouseMoveCallback

	glfwWindow *glfw.Window
}

type RunCallback func(delta float64)
type MousePressCallback func(state glfw.Action, x, y float64)
type MouseMoveCallback func(x, y float64)

func NewWindow(width, height int, title string) (*Window, error) {
	w := Window{}
	w.width = width
	w.height = height
	w.title = title

	err := w.setupGLFW()
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *Window) cursorPress(glfwWindow *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	x, y := glfwWindow.GetCursorPos()
	if w.mousePress != nil {
		w.mousePress(action, x, y)
	}
}

func (w *Window) cursorMove(glfwWindow *glfw.Window, xpos float64, ypos float64) {
	x, y := glfwWindow.GetCursorPos()
	if w.mouseMove != nil {
		w.mouseMove(x, y)
	}
}

func (w *Window) setupGLFW() error {
	if err := glfw.Init(); err != nil {
		return errors.New(fmt.Sprint("failed to initialize glfw:", err))
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create glfw window:", err))
	}

	window.SetMouseButtonCallback(w.cursorPress)
	window.SetCursorPosCallback(w.cursorMove)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return errors.New(fmt.Sprint("failed to initialize OpenGL: ", err))
	}

	w.glfwWindow = window

	return nil

}

func (w *Window) SetRunStep(call RunCallback) {
	w.step = call
}

func (w *Window) SetMousePress(call MousePressCallback) {
	w.mousePress = call
}

func (w *Window) SetMouseMove(call MouseMoveCallback) {
	w.mouseMove = call
}

func (w *Window) Run() {
	w.previousTime = glfw.GetTime()

	for !w.glfwWindow.ShouldClose() { //|| escape
		// Calculate the delta time to give the step callback
		time := glfw.GetTime()
		delta := time - w.previousTime
		w.previousTime = time

		if w.step != nil {
			w.step(delta)
		}

		w.glfwWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func (w *Window) Destroy() {
	glfw.Terminate()
}
