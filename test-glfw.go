// run with: go run test-glfw.go shaders.go

package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	"log"
	"runtime"
	// "github.com/go-gl/mathgl/mgl32"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const windowWidth = 800
const windowHeight = 600

func main() {
	// fmt.Println(vertexShaderT)
	// fmt.Println(fragmentShader)

	mki3dPtr, err := mki3d.ReadFile("noname.mki3d")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", mki3d.Stringify(mki3dPtr))

	// fragments from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) // test quality
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
	gl.Enable(gl.MULTISAMPLE) // test (probably not needed ...)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.3, 1.0)

	// callbacks

	// test Shader
	mki3dShaderPtr, err := tmki3d.MakeShader()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *mki3dShaderPtr)    // test
	fmt.Printf("%+v\n", mki3dShaderPtr.Seg) // test
	fmt.Printf("%+v\n", mki3dShaderPtr.Tr)  // test

	// test GLBuf

	mki3dGLBufPtr, err := tmki3d.MakeGLBuf(mki3dPtr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", *mki3dGLBufPtr) // test

	// test ViewMatrix
	// fmt.Println(ViewMatrix(mki3dPtr.View))

	// test SetFromMki3d
	/// var mki3dGLUni tmki3d.GLUni
	/// mki3dGLUni.SetFromMki3d(mki3dPtr, 100, 100)
	width, height := window.GetSize()
	mki3dGLUniPtr, err := tmki3d.MakeGLUni(mki3dPtr, width, height)
	fmt.Printf("%+v\n", *mki3dGLUniPtr)

	mki3dDataShaderTrPtr, err := tmki3d.MakeDataShaderTr(mki3dShaderPtr.Tr, &(mki3dGLBufPtr.Tr), mki3dGLUniPtr, mki3dPtr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *mki3dDataShaderTrPtr) // test

	SizeCallback := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		fmt.Println(width, height)
		// fmt.Println(tmki3d.ProjectionMatrix(mki3dPtr.Projection, width, height))
		// mki3dGLUni.SetFromMki3d(mki3dPtr, width, height)
		mki3dGLUniPtr.ProjectionUni = tmki3d.ProjectionMatrix(mki3dPtr.Projection, width, height)
		fmt.Printf("%+v\n", *mki3dGLUniPtr)

		mki3dDataShaderTrPtr.DrawStage()

	}

	// setting callbacks
	window.SetSizeCallback(SizeCallback)

	previousTime := glfw.GetTime()
	// main loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		_ = elapsed // do not forget!

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
		mki3dDataShaderTrPtr.DrawStage()

		// Maintenance
		window.SwapBuffers()
		glfw.WaitEvents()
		glfw.PollEvents()
	}
}
