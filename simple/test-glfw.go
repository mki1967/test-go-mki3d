// This program demontsrates the use of packages: "github.com/mki1967/go-mki3d/mki3d" and  "github.com/mki1967/go-mki3d/glmki3d"
// It reads mki3d data from a file  produced with MKI3D RAPID MODELER ( https://mki1967.github.io/mki3d/ ) 
// and displays it using the packages "github.com/go-gl/gl/v3.3-core/gl", "github.com/go-gl/glfw/v3.2/glfw" and "github.com/go-gl/mathgl/mgl32"
// 
// Run in the source code directory with: go run *.go  <filename>.mki3d
// 
package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-mki3d/glmki3d"
	"log"
	"os"
	"runtime"
	// "github.com/go-gl/mathgl/mgl32"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const windowWidth = 800
const windowHeight = 600

var DataShaderPtr *glmki3d.DataShader // global variable in the main package

var Window *glfw.Window // main window

func message(msg string) error {
	fmt.Println(msg)
	err := Window.Iconify()
	if err != nil {
		return err
	}
	fmt.Print("(PRESS ENTER TO RESUME:)")
	fmt.Scanln()
	err = Window.Restore()
	if err != nil {
		panic(err)
	}
	Window.Show()
	fmt.Println("RESUMED.")
	return err
}

var doInMainThread func() = nil

func main() {

	// get file name from command line argument
	if len(os.Args) < 2 {
		panic(errors.New(" *** PROVIDE FILE NAME AS A COMMAND LINE ARGUMENT !!! *** "))
	}
	fmt.Println("Trying to read from ", os.Args[1])

	// Load mki3d data from a file
	mki3dPtr, err := mki3d.ReadFile(os.Args[1]) // load mki3d data from the file produced by MKI3D
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%v\n", mki3d.Stringify(mki3dPtr)) // for tests ... demonstrates  mki3d.Stringify

	// updated fragments from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) // try multisampling for better quality ...
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	Window = window // copy to global variable
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.MULTISAMPLE) // probably not needed ...
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.3, 1.0)

	// Create both (segment and triangle) shaders in single structure
	mki3dShaderPtr, err := glmki3d.MakeShader()
	if err != nil {
		panic(err)
	}

	mki3dDataShaderPtr, err := glmki3d.MakeDataShader(mki3dShaderPtr, mki3dPtr)
	if err != nil {
		panic(err)
	}

	DataShaderPtr = mki3dDataShaderPtr // set the global variable

	// Get current width and height of the window
	width, height := window.GetSize() // needed for projection setting
	DataShaderPtr.UniPtr.SetProjectionFromMki3d(mki3dPtr, width, height) // set the sane projection as in the mki3d data
	DataShaderPtr.UniPtr.SetViewFromMki3d(mki3dPtr) // set the same view as in mki3d data
	DataShaderPtr.UniPtr.SetLightFromMki3d(mki3dPtr) // set the same light as in mki3d data 

	// setting callbacks
	window.SetSizeCallback(SizeCallback)
	window.SetKeyCallback(KeyCallback)

	message(helpText) // initial help message

	previousTime := glfw.GetTime()
	// main loop
	for !window.ShouldClose() {

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		_ = elapsed // do not forget!

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
		mki3dDataShaderPtr.DrawStage()

		// Maintenance
		window.SwapBuffers()
		glfw.WaitEvents()
		// glfw.PollEvents()
	}
}
