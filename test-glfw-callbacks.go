package main

import (
	// "fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	// "github.com/go-gl/mathgl/mgl32"
)

// Function to be used as resize callback
func SizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))                                                                 // inform GL about resize
	DataShaderPtr.UniPtr.ProjectionUni = tmki3d.ProjectionMatrix(DataShaderPtr.Mki3dPtr.Projection, width, height) // recompute projection matrix
}
