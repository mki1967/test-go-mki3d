package main

import (
	//"fmt" // tests
	// "errors"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	"math"
	// "math/rand"
)

var MonsterSpeed float32 = 10

// Parameters of a single monster
type MonsterType struct {
	Position mgl32.Vec3         // current position
	Speed    mgl32.Vec3         // speed vector
	DSPtr    *tmki3d.DataShader // shape for redraw (may be shared by many)
	time     float64            // last update time
}

// Creates a monster wth random speed direction at position pos with datashader *dsptr
func MakeMonster(pos mgl32.Vec3, dsPtr *tmki3d.DataShader) *MonsterType {
	var m MonsterType
	m.Position = pos
	m.DSPtr = dsPtr
	m.Speed = RandRotated(mgl32.Vec3{0, 0, MonsterSpeed})
	m.time = glfw.GetTime()
	return &m
}

// Redraw monster m
func (m *MonsterType) Draw() {
	m.DSPtr.UniPtr.SetModelPosition(m.Position)
	m.DSPtr.DrawModel()
}

// Update monster m in game g
func (m *MonsterType) Update(g *Mki3dGame) {
	now := glfw.GetTime()
	elapsed := float32(now - m.time)
	m.time = now

	dv := m.Speed.Mul(elapsed)
	m.Position = m.Position.Add(dv)
	if m.Position[0] >= g.VMax[0] {
		m.Speed[0] = float32(-math.Abs(float64(m.Speed[0])))
	}
	if m.Position[0] <= g.VMin[0] {
		m.Speed[0] = float32(math.Abs(float64(m.Speed[0])))
	}

	if m.Position[1] >= g.VMax[1] {
		m.Speed[1] = float32(-math.Abs(float64(m.Speed[1])))
	}
	if m.Position[1] <= g.VMin[1] {
		m.Speed[1] = float32(math.Abs(float64(m.Speed[1])))
	}

	if m.Position[2] >= g.VMax[2] {
		m.Speed[2] = float32(-math.Abs(float64(m.Speed[2])))
	}
	if m.Position[2] <= g.VMin[2] {
		m.Speed[2] = float32(math.Abs(float64(m.Speed[2])))
	}
}
