package main

import (
	// "fmt" // tests
	// "errors"
	// "github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/test-go-mki3d/tmki3d"
	"math"
	"math/rand"
)

// RotHVType represents sequence of two rotations:
// by the angle XY on XY-plane and by the angle YZ on YZ-plane
// (in degrees)
type RotHVType struct {
	XZ float64
	YZ float64
}

const degToRadians = math.Pi / 180

func (rot *RotHVType) WorldRotatedVector(vector mgl32.Vec3) mgl32.Vec3 {
	c1 := float32(math.Cos(rot.XZ * degToRadians))
	s1 := float32(math.Sin(rot.XZ * degToRadians))
	c2 := float32(math.Cos(rot.YZ * degToRadians))
	s2 := float32(math.Sin(rot.YZ * degToRadians))

	return mgl32.Vec3{
		c1*vector[0] - s1*s2*vector[1] - s1*c2*vector[2],
		c2*vector[1] - s2*vector[2],
		s1*vector[0] + c1*s2*vector[1] + c1*c2*vector[2],
	}
}

func (rot *RotHVType) ViewerRotatedVector(vector mgl32.Vec3) mgl32.Vec3 {
	c1 := float32(math.Cos(-rot.XZ * degToRadians))
	s1 := float32(math.Sin(-rot.XZ * degToRadians))
	c2 := float32(math.Cos(-rot.YZ * degToRadians))
	s2 := float32(math.Sin(-rot.YZ * degToRadians))

	return mgl32.Vec3{
		c1*vector[0] - s1*vector[2],
		-s2*s1*vector[0] + c2*vector[1] - s2*c1*vector[2],
		c2*s1*vector[0] + s2*vector[1] + c2*c1*vector[2],
	}
}

func RandRotated(vec mgl32.Vec3) mgl32.Vec3 {
	var rot RotHVType
	rot.XZ = rand.Float64() * 360
	rot.YZ = rand.Float64() * 360
	return rot.WorldRotatedVector(vec)
}
