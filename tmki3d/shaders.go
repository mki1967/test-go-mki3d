package tmki3d

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/mki3d"
	"strings"
)

// Vertex shader for drawing triangles
var vertexShaderT = `
#version 330

/* attributes */
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec3 color;

/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
uniform vec3 light;
uniform float ambient; 
 
/* output to fragment shader */
out vec4 vColor;

void main() {
    /* compute shaded color */
    vec4 modelNormal=model*vec4(normal, 1);
    float shade= abs( dot( modelNormal.xyz, light ) ); 
    vColor= (ambient+(1.0-ambient)*shade)*vec4(color, 1.0);
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

// vertex shader for drawing segments
var vertexShaderS = `
#version 330

/* attributes */
layout (location = 0) in vec3 position;
layout (location = 2) in vec3 color;

/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
 
/* output to fragment shader */
out vec4 vColor;

void main() {
    /* compute shaded color */
    vColor= vec4(color, 1.0);
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

// fragment shader - the same for segments and triangles
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

// structure for mki3d shader for drawing triangles
// with references to attributes and uniform locations.
type Mki3dShaderTr struct {
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
	AmbientUni    int32
}

// MakeMki3dShaderTr compiles  mki3d shader and
// returns Mki3dShaderTr structure with reference to the program and its attributes and uniforms
// or error
func MakeMki3dShaderTr() (shaderPtr *Mki3dShaderTr, err error) {
	program, err := newProgram(vertexShaderT, fragmentShader)
	if err != nil {
		return nil, err
	}

	var shader Mki3dShaderTr

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

// TO DO: Mki3dShaderSeg, MakeMki3dShaderSeg() ...
// Mki3dShaderSeg is a structure for mki3d shader for drawing segments
// with references to attributes and uniform locations.
type Mki3dShaderSeg struct {
	// program Id
	ProgramId uint32
	// locations of attributes
	PositionAttr uint32
	ColorAttr    uint32
	// locations of uniforms ( why int32 instead of uint32 ? )
	ProjectionUni int32
	ViewUni       int32
	ModelUni      int32
}

// MakeMki3dShaderSeg compiles  mki3d shader and
// returns Mki3dShaderSeg structure with reference to the program and its attributes and uniforms
// or error
func MakeMki3dShaderSeg() (shaderPtr *Mki3dShaderSeg, err error) {
	program, err := newProgram(vertexShaderS, fragmentShader)
	if err != nil {
		return nil, err
	}

	var shader Mki3dShaderSeg

	// set ProgramId
	shader.ProgramId = program

	// set attributes
	shader.PositionAttr = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	shader.ColorAttr = uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))

	// set uniforms
	shader.ProjectionUni = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	shader.ViewUni = gl.GetUniformLocation(program, gl.Str("view\x00"))
	shader.ModelUni = gl.GetUniformLocation(program, gl.Str("model\x00"))
	return &shader, nil
}

// Both shaders in one struct
type Mki3dShader struct {
	Seg *Mki3dShaderSeg
	Tr  *Mki3dShaderTr
}

func MakeMki3dShader() (shaderPtr *Mki3dShader, err error) {
	shaderSeg, err := MakeMki3dShaderSeg()
	if err != nil {
		return nil, err
	}

	shaderTr, err := MakeMki3dShaderTr()
	if err != nil {
		return nil, err
	}

	return &Mki3dShader{Seg: shaderSeg, Tr: shaderTr}, err

}
