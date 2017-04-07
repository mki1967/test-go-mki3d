package main

import (
	// "fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
)

// Function to be used as resize callback
func SizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))                                                                 // inform GL about resize
	DataShaderPtr.UniPtr.ProjectionUni = tmki3d.ProjectionMatrix(DataShaderPtr.Mki3dPtr.Projection, width, height) // recompute projection matrix
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release {
		return
	}

	const angle = 1 * math.Pi / 180
	const step = 0.5

	switch {

	/* rotate */
	case key == glfw.KeyRight && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.HomogRotate3DY(-angle).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyLeft && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.HomogRotate3DY(angle).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyUp && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.HomogRotate3DX(angle).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyDown && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.HomogRotate3DX(-angle).Mul4(DataShaderPtr.UniPtr.ViewUni)

		/* move */
	case key == glfw.KeyRight && mods == glfw.ModShift:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(-step, 0, 0).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyLeft && mods == glfw.ModShift:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(step, 0, 0).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyUp && mods == glfw.ModShift:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(0, -step, 0).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyDown && mods == glfw.ModShift:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(0, step, 0).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyF && mods == glfw.ModShift:
		fallthrough
	case key == glfw.KeyF && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(0, 0, -step).Mul4(DataShaderPtr.UniPtr.ViewUni)
	case key == glfw.KeyB && mods == glfw.ModShift:
		fallthrough
	case key == glfw.KeyB && mods == 0:
		DataShaderPtr.UniPtr.ViewUni = mgl32.Translate3D(0, 0, step).Mul4(DataShaderPtr.UniPtr.ViewUni)
	}
}
