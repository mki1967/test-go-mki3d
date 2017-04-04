package main

import (
	// "fmt"
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Mki3dDataShaderTr is a binding between data and shader for triangles
type Mki3dDataShaderTr struct {
	ShaderPtr *Mki3dShaderTr // pointer to the GL shader program structure
	VAO       uint32         // GL Vertex Array Object
	BufPtr    *Mki3dGLBufTr  // pointer to GL buffers structure
	UniPtr    *Mki3dGLUni    // pointer to GL uniform parameters structure

}

// MakeMki3dDataShaderTr either returns a pointer to anewly created Mki3dDataShaderTr or an error.
// The parameters should be pointers to existing and initiated objects
func MakeMki3dDataShaderTr(sPtr *Mki3dShaderTr, bPtr *Mki3dGLBufTr, uPtr *Mki3dGLUni) (dsPtr *Mki3dDataShaderTr, err error) {
	if sPtr == nil {
		return nil, errors.New("sPtr == nil // type *Mki3dShaderTr ")
	}
	if bPtr == nil {
		return nil, errors.New("bPtr == nil // type *Mki3dGLBufTr ")
	}
	if uPtr == nil {
		return nil, errors.New("uPtr == nil // type *Mki3dGLUni ")
	}

	ds := Mki3dDataShaderTr{ShaderPtr: sPtr, BufPtr: bPtr, UniPtr: uPtr}
	err = ds.InitVAO()
	if err != nil {
		return nil, err
	}

	// OK
	return &ds, nil

}

// InitVAO init the VAO field of ds. ds, ds.ShaderPtr  and ds.BufPtr must be not nil and previously initiated
func (ds *Mki3dDataShaderTr) InitVAO() (err error) {
	if ds == nil {
		return errors.New("ds == nil // type  *Mki3dDataShaderTr ")
	}

	if ds.BufPtr == nil {
		return errors.New("ds.BufPtr == nil // type *Mki3dGLBufTr")
	}

	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type ")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.GenVertexArrays(1, &(ds.VAO))
	gl.BindVertexArray(ds.VAO)

	/* EXAMPLE:
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	*/

	// bind vertex positions
	gl.BindBuffer(gl.ARRAY_BUFFER, ds.BufPtr.PositionBuf)
	gl.EnableVertexAttribArray(ds.ShaderPtr.PositionAttr)
	gl.VertexAttribPointer(ds.ShaderPtr.PositionAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// bind vertex colors
	gl.BindBuffer(gl.ARRAY_BUFFER, ds.BufPtr.ColorBuf)
	gl.EnableVertexAttribArray(ds.ShaderPtr.ColorAttr)
	gl.VertexAttribPointer(ds.ShaderPtr.ColorAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// bind vertex normals
	gl.BindBuffer(gl.ARRAY_BUFFER, ds.BufPtr.NormalBuf)
	gl.EnableVertexAttribArray(ds.ShaderPtr.NormalAttr)
	gl.VertexAttribPointer(ds.ShaderPtr.NormalAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// ...
	gl.BindVertexArray(0) //
	// TO DO: ...
	return nil

}
