package tmki3d

import (
	// "fmt"
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
)

// DataShaderTr is a binding between data and shader for triangles
type DataShaderTr struct {
	ShaderPtr *ShaderTr // pointer to the GL shader program structure
	VAO       uint32         // GL Vertex Array Object
	BufPtr    *GLBufTr  // pointer to GL buffers structure
	UniPtr    *GLUni    // pointer to GL uniform parameters structure

}

// MakeDataShaderTr either returns a pointer to anewly created DataShaderTr or an error.
// The parameters should be pointers to existing and initiated objects
// MakeDataShaderTr inits its VAO
func MakeDataShaderTr(sPtr *ShaderTr, bPtr *GLBufTr, uPtr *GLUni) (dsPtr *DataShaderTr, err error) {
	if sPtr == nil {
		return nil, errors.New("sPtr == nil // type *ShaderTr ")
	}
	if bPtr == nil {
		return nil, errors.New("bPtr == nil // type *GLBufTr ")
	}
	if uPtr == nil {
		return nil, errors.New("uPtr == nil // type *GLUni ")
	}

	ds := DataShaderTr{ShaderPtr: sPtr, BufPtr: bPtr, UniPtr: uPtr}
	err = ds.InitVAO()
	if err != nil {
		return nil, err
	}

	ds.InitVAO()
	ds.LightToShader()
	// OK
	return &ds, nil

}

// LightUniToShader sets  light uniform parameters from ds.UniPtr to ds.ShaderPtr  (both must be not nil and previously initiated)
func (ds *DataShaderTr) LightToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.Uniform3fv(ds.ShaderPtr.LightUni, 1, &(ds.UniPtr.LightUni[0]))
	gl.Uniform1f(ds.ShaderPtr.AmbientUni, ds.UniPtr.AmbientUni)

	return nil
}

// InitVAO init the VAO field of ds. ds, ds.ShaderPtr  and ds.BufPtr must be not nil and previously initiated
func (ds *DataShaderTr) InitVAO() (err error) {
	if ds == nil {
		return errors.New("ds == nil // type  *DataShaderTr ")
	}

	if ds.BufPtr == nil {
		return errors.New("ds.BufPtr == nil // type *GLBufTr")
	}

	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
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
