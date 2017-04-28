package main

import (
	// "fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
)

// Function to be used as resize callback
func SizeCallback(w *glfw.Window, width int, height int) {
	g := GamePtr                                                                                                 // short name
	gl.Viewport(0, 0, int32(width), int32(height))                                                               // inform GL about resize
	g.StageDSPtr.UniPtr.ProjectionUni = tmki3d.ProjectionMatrix(g.StageDSPtr.Mki3dPtr.Projection, width, height) // recompute projection matrix
	// fmt.Println("SizeCallback ",  width, " ", height)
}

