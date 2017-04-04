package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
	// "strings"
)

// references to the objects defining the shape and parameters of mki3d object

// Mki3dGLBufTr contains references to GL triangle buffers for triangle shader's input attributes
type Mki3dGLBufTr struct {
	// buffer objects in GL
	// triangles:
	VertexCount int32 // the last argument for gl.DrawArrays
	PositionBuf uint32
	NormalBuf   uint32
	ColorBuf    uint32
}



// Mki3dGLBufSeg contains references to GL segment buffers for segment shader's input attributes
type Mki3dGLBufSeg struct {
	// buffer objects in GL
	// segments:
	VertexCount int32 // the last argument for gl.DrawArrays
	PositionBuf uint32
	ColorBuf    uint32
}

// Mki3dGLBuf contains references to GL buffers for shaders' input attributes
type Mki3dGLBuf struct {
	// buffer objects in GL
	// triangles:
	Triangles Mki3dGLBufTr
	// segments:
	Segments Mki3dGLBufSeg
}


func (glBuf *Mki3dGLBufTr) LoadTriangleBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.VertexCount = int32(3 * len(mki3dData.Model.Triangles))
	dataPos := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataCol := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataNor := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	i := 0
	for _, triangle := range mki3dData.Model.Triangles {
		// compute normal
		a := mgl32.Vec3(triangle[0].Position)
		b := mgl32.Vec3(triangle[1].Position)
		c := mgl32.Vec3(triangle[2].Position)
		normal := (b.Sub(a)).Cross(c.Sub(a))
		if normal.Dot(normal) > 0 {
			normal = normal.Normalize()
		}
		// fmt.Println( "normal: ", normal ) /// test ...
		// append to buffers
		for j := 0; j < 3; j++ {
			dataPos = append(dataPos, triangle[j].Position[0:3]...)
			dataCol = append(dataCol, triangle[j].Color[0:3]...)
			dataNor = append(dataNor, normal[0:3]...)
			i = i + 3
		}
	}

	fmt.Println(dataPos) // test
	fmt.Println(dataCol) // test
	fmt.Println(dataNor) // test

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.PositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.ColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.NormalBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataNor)*4 /* 4 bytes per flat32 */, gl.Ptr(dataNor), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

func (glBuf *Mki3dGLBufSeg) LoadSegmentBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.VertexCount = int32(2 * len(mki3dData.Model.Segments))
	dataPos := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	dataCol := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	i := 0
	for _, segment := range mki3dData.Model.Segments {
		for j := 0; j < 2; j++ {
			dataPos = append(dataPos, segment[j].Position[0:3]...)
			dataCol = append(dataCol, segment[j].Color[0:3]...)
			i = i + 2
		}
	}

	fmt.Println(dataPos) // test
	fmt.Println(dataCol) // test
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.PositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.ColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

func MakeMki3dGLBufTr(mki3dData *mki3d.Mki3dType) (glBufPtr *Mki3dGLBufTr, err error) {
	var glBuf Mki3dGLBufTr
	var vbo [3]uint32 // 5 is the number of buffers
	gl.GenBuffers(3, &vbo[0])
	// TO DO: test for error ...

	// assign buffer ids from vbo array
	glBuf.PositionBuf = vbo[0]
	glBuf.NormalBuf = vbo[1]
	glBuf.ColorBuf = vbo[2]

	// load data from mki3dData
	glBuf.LoadTriangleBufs(mki3dData)

	return &glBuf, nil
}


func MakeMki3dGLBufSeg(mki3dData *mki3d.Mki3dType) (glBufPtr *Mki3dGLBufSeg, err error) {
	var glBuf Mki3dGLBufSeg
	var vbo [2]uint32 // 5 is the number of buffers
	gl.GenBuffers(2, &vbo[0])
	// TO DO: test for error ...

	// assign buffer ids from vbo array
	glBuf.PositionBuf = vbo[0]
	glBuf.ColorBuf = vbo[1]

	// load data from mki3dData
	glBuf.LoadSegmentBufs(mki3dData)

	return &glBuf, nil
}

func MakeMki3dGLBuf(mki3dData *mki3d.Mki3dType) (glBufPtr *Mki3dGLBuf, err error) {

	glSegBufPtr, err := MakeMki3dGLBufSeg( mki3dData )
	if err != nil {
		return nil, err
	}
	glTrBufPtr, err := MakeMki3dGLBufTr( mki3dData )
	if err != nil {
		return nil, err
	}

	glBuf := Mki3dGLBuf{ Triangles: *glTrBufPtr, Segments:  *glSegBufPtr }
	return &glBuf, nil
}

