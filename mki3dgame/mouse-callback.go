package main

import (
	// "fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/test-go-mki3d/tmki3d"
)

func Mki3dMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Release {
		GamePtr.CancelAction()
		return
	}

	if action == glfw.Press {
		width, height := w.GetSize()
		fx, fy := w.GetCursorPos()
		x := int(fx)
		y := int(fy)
		GamePtr.SetSectorAction(x, y, width, height)
		return
	}

}
