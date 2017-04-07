package tmki3d

import (
	"fmt" // tests
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/mki1967/go-mki3d/mki3d"
)

// DataShaderTr is a binding between data and a shader for triangles
type DataShaderTr struct {
	ShaderPtr *ShaderTr // pointer to the GL shader program structure
	VAO       uint32         // GL Vertex Array Object
	BufPtr    *GLBufTr  // pointer to GL buffers structure
	UniPtr    *GLUni    // pointer to GL uniform parameters structure
	Mki3dPtr    *mki3d.Mki3dType    // pointer to original Mki3dType data

}

// DataShaderSeg is a binding between data and a shader for triangles
type DataShaderSeg struct {
	ShaderPtr *ShaderSeg // pointer to the GL shader program structure
	VAO       uint32         // GL Vertex Array Object
	BufPtr    *GLBufSeg  // pointer to GL buffers structure
	UniPtr    *GLUni    // pointer to GL uniform parameters structure
	Mki3dPtr    *mki3d.Mki3dType    // pointer to original Mki3dType data

}

// MakeDataShaderTr either returns a pointer to anewly created DataShaderTr or an error.
// The parameters should be pointers to existing and initiated objects
// MakeDataShaderTr inits its VAO
func MakeDataShaderTr(sPtr *ShaderTr, bPtr *GLBufTr, uPtr *GLUni, mPtr *mki3d.Mki3dType ) (dsPtr *DataShaderTr, err error) {
	if sPtr == nil {
		return nil, errors.New("sPtr == nil // type *ShaderTr ")
	}
	if bPtr == nil {
		return nil, errors.New("bPtr == nil // type *GLBufTr ")
	}
	if uPtr == nil {
		return nil, errors.New("uPtr == nil // type *GLUni ")
	}

	if mPtr == nil {
		return nil, errors.New("mPtr == nil // type *Mki3dType ")
	}

	ds := DataShaderTr{ShaderPtr: sPtr, BufPtr: bPtr, UniPtr: uPtr, Mki3dPtr: mPtr  }
	err = ds.InitVAO()
	if err != nil {
		return nil, err
	}

	ds.InitVAO()

	return &ds, nil

}

// MakeDataShaderSeg either returns a pointer to anewly created DataShaderTr or an error.
// The parameters should be pointers to existing and initiated objects
// MakeDataShaderTr inits its VAO
func MakeDataShaderSeg(sPtr *ShaderSeg, bPtr *GLBufSeg, uPtr *GLUni, mPtr *mki3d.Mki3dType ) (dsPtr *DataShaderSeg, err error) {
	if sPtr == nil {
		return nil, errors.New("sPtr == nil // type *ShaderTr ")
	}
	if bPtr == nil {
		return nil, errors.New("bPtr == nil // type *GLBufTr ")
	}
	if uPtr == nil {
		return nil, errors.New("uPtr == nil // type *GLUni ")
	}

	if mPtr == nil {
		return nil, errors.New("mPtr == nil // type *Mki3dType ")
	}

	ds := DataShaderSeg{ShaderPtr: sPtr, BufPtr: bPtr, UniPtr: uPtr, Mki3dPtr: mPtr  }
	err = ds.InitVAO()
	if err != nil {
		return nil, err
	}

	ds.InitVAO()

	return &ds, nil

}



// UniLightToShader sets  light uniform parameters from ds.UniPtr to ds.ShaderPtr  (both must be not nil and previously initiated)
func (ds *DataShaderTr) UniLightToShader() (err error) {
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


// UniModelToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderTr) UniModelToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ModelUni, 1, false, &(ds.UniPtr.ModelUni[0]) )

	return nil
}

// UniModelToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderSeg) UniModelToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ModelUni, 1, false, &(ds.UniPtr.ModelUni[0]) )

	return nil
}

// UniViewToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderTr) UniViewToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ViewUni, 1, false, &(ds.UniPtr.ViewUni[0]) )

	return nil
}

// UniViewToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderSeg) UniViewToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ViewUni, 1, false, &(ds.UniPtr.ViewUni[0]) )

	return nil
}

// UniProjectionToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr 
func (ds *DataShaderTr) UniProjectionToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ProjectionUni, 1, false, &(ds.UniPtr.ProjectionUni[0]) )

	return nil
}

// UniProjectionToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr 
func (ds *DataShaderSeg) UniProjectionToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ProjectionUni, 1, false, &(ds.UniPtr.ProjectionUni[0]) )

	return nil
}

// InitStage initiates stage parameters in ds.ShaderPtr assuming that ds is a stage
func (ds *DataShaderTr) InitStage() (err error) {
	if ds.Mki3dPtr == nil {
		return errors.New("ds.Mki3dPtr == nil // type *Mki3dType")
	}

	err = ds.UniProjectionToShader() // set projection
	if err != nil {
		return err
	}
	
	err = ds.UniViewToShader() // set view
	if err != nil {
		return err
	}
	
	err = ds.UniLightToShader() // set light - for triangles only
	if err != nil {
		return err
	}

	// to be moved elsewhere ...
	bg := ds.Mki3dPtr.BackgroundColor
	gl.ClearColor(bg[0], bg[1], bg[2], 1.0)

	return nil
	
}

// InitStage initiates stage parameters in ds.ShaderPtr assuming that ds is a stage
func (ds *DataShaderSeg) InitStage() (err error) {
	if ds.Mki3dPtr == nil {
		return errors.New("ds.Mki3dPtr == nil // type *Mki3dType")
	}

	err = ds.UniProjectionToShader() // set projection
	if err != nil {
		return err
	}
	
	err = ds.UniViewToShader() // set view
	if err != nil {
		return err
	}

	// to be moved elsewhere ...
	bg := ds.Mki3dPtr.BackgroundColor       
	gl.ClearColor(bg[0], bg[1], bg[2], 1.0)

	return nil
	
}

func (ds *DataShaderTr) DrawModel() {
	if ds.BufPtr.VertexCount == 0 {
		return // nothing to draw
	}
	ds.UniModelToShader()
	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.BindVertexArray(ds.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, ds.BufPtr.VertexCount)
	fmt.Println( "ds.BufPtr.VertexCount: ", ds.BufPtr.VertexCount)
	fmt.Println( "ds.VAO: ", ds.VAO)
	gl.BindVertexArray(0)
}

func (ds *DataShaderSeg) DrawModel() {
	if ds.BufPtr.VertexCount == 0 {
		return // nothing to draw
	}
	ds.UniModelToShader()
	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.BindVertexArray(ds.VAO)
	gl.DrawArrays(gl.LINES, 0, ds.BufPtr.VertexCount)
	fmt.Println( "ds.BufPtr.VertexCount: ", ds.BufPtr.VertexCount)
	fmt.Println( "ds.VAO: ", ds.VAO)
	gl.BindVertexArray(0)
}

func (ds *DataShaderTr) DrawStage() {
	ds.InitStage()
	ds.DrawModel()
}

func (ds *DataShaderSeg) DrawStage() {
	ds.InitStage()
	ds.DrawModel()
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

	gl.BindVertexArray(0) // unbind VAO

	return nil

}

// InitVAO init the VAO field of ds. ds, ds.ShaderPtr  and ds.BufPtr must be not nil and previously initiated
func (ds *DataShaderSeg) InitVAO() (err error) {
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


	// bind vertex positions
	gl.BindBuffer(gl.ARRAY_BUFFER, ds.BufPtr.PositionBuf)
	gl.EnableVertexAttribArray(ds.ShaderPtr.PositionAttr)
	gl.VertexAttribPointer(ds.ShaderPtr.PositionAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// bind vertex colors
	gl.BindBuffer(gl.ARRAY_BUFFER, ds.BufPtr.ColorBuf)
	gl.EnableVertexAttribArray(ds.ShaderPtr.ColorAttr)
	gl.VertexAttribPointer(ds.ShaderPtr.ColorAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	gl.BindVertexArray(0) // unbind VAO

	return nil

}
