package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/mki1967/go-mki3d/mki3d"
	"strings"
)

var vertexShader = `
#version 330

/* attributes */
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec3 color;

/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
uniform vec4 light;
uniform float ambient; 
 
/* output to fragment shader */
out vec4 vColor;

void main() {
    /* compute shaded color */
    vec4 modelNormal=model*vec4(normal, 1);
    float shade= abs( dot( modelNormal, light ) ); 
    vColor= (ambient+(1.0-ambient)*shade)*vec4(color, 1.0);
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

var fragmentShader = `
#version 330

/* input from vertex shader */
in vec4 vColor;

/* fragment color output */
out vec4 outputColor;

void main() {
    outputColor = vColor ;
}
` + "\x00"

// from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// structure for standard mki3d shader with references to attributes and uniform locations
type Mki3dShader struct {
	// program Id
	ProgramId uint32
	// locations of attributes
	PositionAttr uint32
	NormalAttr   uint32
	ColorAttr    uint32
	// locations of uniforms ( why int32 instead of uint32 ? )
	ProjectionUni int32
	ViewUni       int32
	ModelUni      int32
	LightUni      int32
}

// MakeMki3dShader compiles standard mki3d shader and
// returns Mki3dShader structure with reference to the program and its attributes and uniforms
// or error
func MakeMki3dShader() (shaderPtr *Mki3dShader, err error) {
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	var shader Mki3dShader

	// set ProgramId
	shader.ProgramId = program

	// set attributes
	shader.PositionAttr = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	shader.NormalAttr = uint32(gl.GetAttribLocation(program, gl.Str("normal\x00")))
	shader.ColorAttr = uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))

	// set uniforms
	shader.ProjectionUni = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	shader.ViewUni = gl.GetUniformLocation(program, gl.Str("view\x00"))
	shader.ModelUni = gl.GetUniformLocation(program, gl.Str("model\x00"))
	shader.LightUni = gl.GetUniformLocation(program, gl.Str("light\x00"))
	return &shader, nil
}

// references to the objects defining the shape and parameters of mki3d object
type Mki3dGLBuf struct {
	// buffer objects in GL
	// triangles:
	trianglePositionBuf uint32
	triangleNormalBuf   uint32
	triangleColorBuf    uint32
	// segments:
	segmentPositionBuf uint32
	segmentColorBuf    uint32
}

func (glBuf *Mki3dGLBuf) LoadTrianglePositions(mki3dData *mki3d.Mki3dType) {
	data := make([]float32, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.trianglePositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4 /* 4 bytes per flat32 */, gl.Ptr(data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

func MakeMki3dGLBuf(mki3dData *mki3d.Mki3dType) (glBufPtr *Mki3dGLBuf, err error) {
	var glBuf Mki3dGLBuf
	var vbo [5]uint32 // 5 is the number of buffers
	gl.GenBuffers(5, &vbo[0])

	// assign buffer ids from vbo array
	glBuf.trianglePositionBuf = vbo[0]
	glBuf.triangleNormalBuf = vbo[1]
	glBuf.triangleColorBuf = vbo[2]
	glBuf.segmentPositionBuf = vbo[3]
	glBuf.segmentColorBuf = vbo[4]

	// load data from mki3dData
	glBuf.LoadTrianglePositions(mki3dData)
	// ...

	return &glBuf, nil
}
