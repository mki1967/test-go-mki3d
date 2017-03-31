package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
	// "strings"
)

// references to the objects defining the shape and parameters of mki3d object

// Mki3dGLBuf contains references to GL buffers for shaders' input attributes
type Mki3dGLBuf struct {
	// buffer objects in GL
	// triangles:
	triangleVertexCount int32 // the last argument for gl.DrawArrays
	trianglePositionBuf uint32
	triangleNormalBuf   uint32
	triangleColorBuf    uint32
	// segments:
	segmentVertexCount int32 // the last argument for gl.DrawArrays
	segmentPositionBuf uint32
	segmentColorBuf    uint32
}

func (glBuf *Mki3dGLBuf) LoadTriangleBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.triangleVertexCount = int32(3 * len(mki3dData.Model.Triangles))
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

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.trianglePositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.triangleColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.triangleNormalBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataNor)*4 /* 4 bytes per flat32 */, gl.Ptr(dataNor), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

func (glBuf *Mki3dGLBuf) LoadSegmentBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.segmentVertexCount = int32(2 * len(mki3dData.Model.Segments))
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
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.segmentPositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.segmentColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)
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
	glBuf.LoadTriangleBufs(mki3dData)
	glBuf.LoadSegmentBufs(mki3dData)
	// ...

	return &glBuf, nil
}

// Mki3dGLUni - values of parameters to be stored in shaders' uniforms
type Mki3dGLUni struct {
	ProjectionUni mgl32.Mat4
	ViewUni       mgl32.Mat4
	ModelUni      mgl32.Mat4
	LightUni      mgl32.Vec3
}

func ProjectionMatrix(p mki3d.ProjectionType, width, height int) mgl32.Mat4 {
	// make both width and height greater than zero
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	h := float32(height)
	w := float32(width)
	xx := p.ZoomY * h / w
	yy := p.ZoomY
	zz := (p.ZFar + p.ZNear) / (p.ZFar - p.ZNear)
	zw := float32(1.0)
	wz := -2 * p.ZFar * p.ZNear / (p.ZFar - p.ZNear)

	var m mgl32.Mat4

	m.SetRow(0, mgl32.Vec4{xx, 0, 0, 0})
	m.SetRow(1, mgl32.Vec4{0, yy, 0, 0})
	m.SetRow(2, mgl32.Vec4{0, 0, zz, wz})
	m.SetRow(3, mgl32.Vec4{0, 0, zw, 0})

	// ...

	return m

}

func MakeMki3dGLUni(mki3dData *mki3d.Mki3dType) (glUniPtr *Mki3dGLUni, err error) {
	var glUni Mki3dGLUni
	glUni.LightUni = mgl32.Vec3(mki3dData.Light.Vector)
	glUni.ModelUni = mgl32.Ident4()
	glUni.ProjectionUni = mgl32.Ident4()

	// ...
	return &glUni, nil
}
