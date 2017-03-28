// run with: go run test-glfw.go shaders.go

package main

import (
	"fmt"
	"github.com/mki1967/go-mki3d/mki3d"
	"runtime"
	"log"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/go-gl/mathgl/mgl32"

)


func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const windowWidth = 800
const windowHeight = 600

func main() {
	fmt.Println(vertexShader)
	fmt.Println(fragmentShader)
	

	mki3dData, err := mki3d.ReadFile("noname.mki3d")
	if err!=nil {
		panic( err )
	}
        fmt.Println(mki3d.Stringify(mki3dData))

	// fragments from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
	
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)


	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.3, 1.0)

	previousTime := glfw.GetTime()
	// main loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		_ =elapsed // do not forget!
		
		/* 
		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
                */
		
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
