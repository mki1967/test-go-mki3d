// run in the source code directory with: go run *.go  <filename>.mki3d

package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix()) // init random generator
}

const windowWidth = 800
const windowHeight = 600

var DataShaderPtr *tmki3d.DataShader // global variable in the main package
var GamePtr *Mki3dGame               // global variable in the main package

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
	// err = Window.Maximize()
	Window.Show()
	fmt.Println("RESUMED.")
	return err
}

var doInMainThread func() = nil

func main() {

	// get file name from command line argument
	if len(os.Args) < 2 {
		panic(errors.New(" *** PROVIDE PATH TO ASSETS DIRECTORY AS A COMMAND LINE ARGUMENT !!! *** "))
	}

	// fmt.Printf("%v\n", mki3d.Stringify(mki3dPtr)) // for tests ...

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

	// callbacks

	// Create both (segment and triangle) shaders in single structure
	mki3dShaderPtr, err := tmki3d.MakeShader()
	if err != nil {
		panic(err)
	}

	// trying to load assets
	fmt.Println("Trying to read assets from ", os.Args[1])

	assetsPtr, err := LoadAssets(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Get current width and height of the window for MakeDataShader
	width, height := window.GetSize()

	game, err := MakeEmptyGame(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = game.Init(width, height)
	if err != nil {
		panic(err)
	}

	GamePtr = game
	// DataShaderPtr = mki3dDataShaderPtr // set the global variable
	DataShaderPtr = game.StageDSPtr // set the global variable

	tokenPtr, err := assetsPtr.LoadRandomToken()
	tokenDataShaderPtr, err := tmki3d.MakeDataShader(mki3dShaderPtr, tokenPtr)

	sectorsPtr, err := assetsPtr.LoadRandomSectors()
	sectorsDataShaderPtr, err := tmki3d.MakeDataShader(mki3dShaderPtr, sectorsPtr)

	sectorsDataShaderPtr.UniPtr.SetSimple()

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
		game.Update()
		game.Redraw()

		// Maintenance
		window.SwapBuffers()
		// glfw.WaitEvents()
		if doInMainThread != nil {
			doInMainThread()     // execute required function
			doInMainThread = nil // done
		}
		glfw.PollEvents()
	}

	// cleanup
	game.StageDSPtr.DeleteData()
	sectorsDataShaderPtr.DeleteData()
	tokenDataShaderPtr.DeleteData()
}
