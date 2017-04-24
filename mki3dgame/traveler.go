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
	// "math/rand"
)

type Traveler struct {
	Position mgl32.Vec3 // position
	Rot      RotHVType  // orientation
}

func (t *Traveler) ViewMatrix() mgl32.Mat4 {
	c1 := float32(math.Cos(-t.Rot.XZ * degToRadians))
	s1 := float32(math.Sin(-t.Rot.XZ * degToRadians))

	c2 := float32(math.Cos(-t.Rot.YZ * degToRadians))
	s2 := float32(math.Sin(-t.Rot.YZ * degToRadians))

	v := t.Rot.ViewerRotatedVector(t.Position.Mul(-1))

	// row-major ??
	return mgl32.Mat4{
		c1, 0, -s1, v[0],
		-s2 * s1, c2, -s2 * c1, v[1],
		c2 * s1, s2, c2 * c1, v[2],
		0, 0, 0, 1,
	}.Transpose()
}

func (t *Traveler) Move(dx, dy, dz float32) {
	v := t.Rot.WorldRotatedVector(mgl32.Vec3{dx, dy, dz})
	t.Position = t.Position.Add(v)

	// check bounds and other conditions ...

}

func MakeTraveler(position mgl32.Vec3) *Traveler {
	var t Traveler
	t.Position = position
	return &t
}
