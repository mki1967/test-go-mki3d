package main

import (
	"fmt" // tests
	// "errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	"math"
	"math/rand"
)

const BoxMargin = 30 // margin for bounding box of the stage

var FrameColor = mki3d.Vector3dType{1.0, 1.0, 1.0} // color of the bounding box frame

var NumberOfMonsters = 6

var NumberOfTokens = 10

// data structure for the game
type Mki3dGame struct {
	// assets info
	AssetsPtr *Assets
	// GL shaders
	ShaderPtr *tmki3d.Shader
	// Shape data shaders
	StageDSPtr   *tmki3d.DataShader
	FrameDSPtr   *tmki3d.DataShader // frame of the bounding box (computed for the stage)
	SectorsDSPtr *tmki3d.DataShader
	TokenDSPtr   *tmki3d.DataShader
	MonsterDSPtr *tmki3d.DataShader

	VMin, VMax mgl32.Vec3 // corners of the bounding box of the stage (computed with the BoxMargin)

	TravelerPtr *Traveler // the first person (the player)

	Monsters []*MonsterType // set of monsters

	Tokens []*TokenType // set of tokens

	TokensRemaining int // number of remaining tokens
}

// Random position in the game stage box with the margin offset from the borders
func (game *Mki3dGame) RandPosition(margin float32) mgl32.Vec3 {
	m := mgl32.Vec3{margin, margin, margin}
	v1 := game.VMin.Add(m)
	v2 := game.VMax.Sub(m)
	return RandPosition(v1, v2)
}

// Random position in the box [vmin, vmax]
func RandPosition(vmin, vmax mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		rand.Float32()*(vmax[0]-vmin[0]) + vmin[0],
		rand.Float32()*(vmax[1]-vmin[1]) + vmin[1],
		rand.Float32()*(vmax[2]-vmin[2]) + vmin[2],
	}
}

// Make game structure with the shader and without any data.
// Prepare assets info using pathToAssets.
// Return pointer to the strucure.
func MakeEmptyGame(pathToAssets string) (*Mki3dGame, error) {
	var game Mki3dGame

	shaderPtr, err := tmki3d.MakeShader()
	if err != nil {
		return nil, err
	}

	game.ShaderPtr = shaderPtr

	assetsPtr, err := LoadAssets(pathToAssets)
	if err != nil {
		return nil, err
	}
	game.AssetsPtr = assetsPtr
	return &game, nil
}

// Load data and init game for the first stage
func (game *Mki3dGame) Init(width, height int) (err error) {
	err = game.InitSectors()
	if err != nil {
		return err
	}

	err = game.InitStage(width, height)
	if err != nil {
		return err
	}

	err = game.InitToken()
	if err != nil {
		return err
	}

	err = game.InitMonster()
	if err != nil {
		return err
	}
	return nil
}

// Load sectors shape and init the SectorsDSPtr.
func (game *Mki3dGame) InitSectors() error {

	sectorsPtr, err := game.AssetsPtr.LoadRandomSectors()
	if err != nil {
		return err
	}

	sectorsDataShaderPtr, err := tmki3d.MakeDataShader(game.ShaderPtr, sectorsPtr)
	if err != nil {
		return err
	}

	sectorsDataShaderPtr.UniPtr.SetSimple()

	if game.SectorsDSPtr != nil {
		game.SectorsDSPtr.DeleteData() // free old GL buffers
	}

	game.SectorsDSPtr = sectorsDataShaderPtr

	return nil
}

// Load token shape and init the tokenDSPtr.
func (game *Mki3dGame) InitToken() error {

	tokenPtr, err := game.AssetsPtr.LoadRandomToken()
	if err != nil {
		return err
	}

	tokenDataShaderPtr, err := tmki3d.MakeDataShader(game.ShaderPtr, tokenPtr)
	if err != nil {
		return err
	}

	tokenDataShaderPtr.UniPtr.SetSimple()

	if game.TokenDSPtr != nil {
		game.TokenDSPtr.DeleteData() // free old GL buffers
	}

	game.TokenDSPtr = tokenDataShaderPtr
	game.GenerateTokens()

	return nil
}

func (game *Mki3dGame) GenerateTokens() {
	game.Tokens = make([]*TokenType, NumberOfTokens)
	for i := range game.Tokens {
		game.Tokens[i] = MakeToken(game.RandPosition(BoxMargin), game.TokenDSPtr)
	}
	game.TokensRemaining = NumberOfTokens
}

func (game *Mki3dGame) DrawTokens() {
	for _, t := range game.Tokens {
		t.Draw()
	}
}

func (game *Mki3dGame) UpdateTokens() {
	for _, t := range game.Tokens {
		t.Update(game)
	}
}

// Load monster shape and init the monsters.
func (game *Mki3dGame) InitMonster() error {

	monsterPtr, err := game.AssetsPtr.LoadRandomMonster()
	if err != nil {
		return err
	}

	monsterDataShaderPtr, err := tmki3d.MakeDataShader(game.ShaderPtr, monsterPtr)
	if err != nil {
		return err
	}

	monsterDataShaderPtr.UniPtr.SetSimple()

	if game.MonsterDSPtr != nil {
		game.MonsterDSPtr.DeleteData() // free old GL buffers
	}

	game.MonsterDSPtr = monsterDataShaderPtr

	// game.MonsterDSPtr.UniPtr.SetModelPosition(game.RandPosition(BoxMargin)) // test
	game.GenerateMonsters()

	return nil
}

func (game *Mki3dGame) GenerateMonsters() {
	game.Monsters = make([]*MonsterType, NumberOfMonsters)
	for i := range game.Monsters {
		game.Monsters[i] = MakeMonster(game.RandPosition(0), game.MonsterDSPtr)
	}
}

func (game *Mki3dGame) DrawMonsters() {
	for _, m := range game.Monsters {
		m.Draw()
	}
}

func (game *Mki3dGame) UpdateMonsters() {
	for _, m := range game.Monsters {
		m.Update(game)
	}
}

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

// Parameters of a single token
type TokenType struct {
	Position  mgl32.Vec3
	Collected bool
	DSPtr     *tmki3d.DataShader // shape for redraw (may be shared by many)
}

// Creates a token  at position pos with datashader *dsptr
func MakeToken(pos mgl32.Vec3, dsPtr *tmki3d.DataShader) *TokenType {
	var t TokenType
	t.Position = pos
	t.DSPtr = dsPtr
	t.Collected = false
	return &t
}

// Redraw token m
func (t *TokenType) Draw() {
	if t.Collected {
		return
	}

	t.DSPtr.UniPtr.SetModelPosition(t.Position)
	t.DSPtr.DrawModel()
}

// square of the distance to collect token
const TokenCollectionSqrDist = 1

// Update monster m in game g
func (t *TokenType) Update(g *Mki3dGame) {
	v := t.Position.Sub(g.TravelerPtr.Position)
	if v.Dot(v) < TokenCollectionSqrDist {
		t.Collected = true
		g.TokensRemaining--
		// some celebrations ...
	}
}

// Load stage shape and init the related data.
func (game *Mki3dGame) InitStage(width, height int) error {

	stagePtr, err := game.AssetsPtr.LoadRandomStage()
	if err != nil {
		return err
	}

	stageDataShaderPtr, err := tmki3d.MakeDataShader(game.ShaderPtr, stagePtr)
	if err != nil {
		return err
	}

	stageDataShaderPtr.UniPtr.SetSimple()
	stageDataShaderPtr.UniPtr.SetProjectionFromMki3d(stagePtr, width, height)
	stageDataShaderPtr.UniPtr.SetLightFromMki3d(stagePtr)

	stageDataShaderPtr.UniPtr.ViewUni = mgl32.Ident4()
	stageDataShaderPtr.UniPtr.ViewUni.SetCol(3, mgl32.Vec3(stageDataShaderPtr.Mki3dPtr.Cursor.Position).Mul(-1).Vec4(1))

	if game.StageDSPtr != nil {
		game.StageDSPtr.DeleteData() // free old GL buffers
	}

	game.StageDSPtr = stageDataShaderPtr

	game.copmuteVMinVMax() // compute bounding box of the stage: VMin, VMax
	game.copmuteFrame()    // visible line frame of the bounding box

	game.TravelerPtr = MakeTraveler(mgl32.Vec3(stagePtr.Cursor.Position))

	return nil
}

// recompute bounding box  with the BoxMargin  corners of the stage.
func (game *Mki3dGame) copmuteVMinVMax() {
	stagePtr := game.StageDSPtr.Mki3dPtr
	game.VMax = mgl32.Vec3(stagePtr.Cursor.Position) // cursror position should be included - the starting poin of traveler
	game.VMin = game.VMax

	for _, seg := range stagePtr.Model.Segments {
		for _, point := range seg {
			for d := range point.Position {
				if game.VMax[d] < point.Position[d] {
					game.VMax[d] = point.Position[d]
				}
				if game.VMin[d] > point.Position[d] {
					game.VMin[d] = point.Position[d]
				}
			}

		}
	}

	for _, tr := range stagePtr.Model.Triangles {
		for _, point := range tr {
			for d := range point.Position {
				if game.VMax[d] < point.Position[d] {
					game.VMax[d] = point.Position[d]
				}
				if game.VMin[d] > point.Position[d] {
					game.VMin[d] = point.Position[d]
				}
			}

		}
	}

	m := mgl32.Vec3{BoxMargin, BoxMargin, BoxMargin}

	game.VMin = game.VMin.Sub(m)
	game.VMax = game.VMax.Add(m)
	fmt.Println(game.VMin, game.VMax) // test
}

// recompute frame of the bounding box corners of the stage.
func (game *Mki3dGame) copmuteFrame() {
	a := game.VMin
	b := game.VMax

	v000 := mki3d.Vector3dType(a)
	v001 := mki3d.Vector3dType{a[0], a[1], b[2]}
	v010 := mki3d.Vector3dType{a[0], b[1], a[2]}
	v011 := mki3d.Vector3dType{a[0], b[1], b[2]}
	v100 := mki3d.Vector3dType{b[0], a[1], a[2]}
	v101 := mki3d.Vector3dType{b[0], a[1], b[2]}
	v110 := mki3d.Vector3dType{b[0], b[1], a[2]}
	v111 := mki3d.Vector3dType(b)

	lines := [][2]mki3d.Vector3dType{
		{v000, v001},
		{v010, v011},
		{v100, v101},
		{v110, v111},

		{v000, v010},
		{v001, v011},
		{v100, v110},
		{v101, v111},

		{v000, v100},
		{v001, v101},
		{v010, v110},
		{v011, v111}}

	segments := mki3d.SegmentsType(make([]mki3d.SegmentType, 12))

	for i := range segments {
		segments[i] = mki3d.SegmentType{
			{Position: lines[i][0], Color: FrameColor},
			{Position: lines[i][1], Color: FrameColor}}
	}

	var frameMki3d mki3d.Mki3dType

	frameMki3d.Model.Segments = segments

	dsPtr, err := tmki3d.MakeDataShader(game.ShaderPtr, &frameMki3d)

	if err != nil {
		panic(err)
	}

	dsPtr.UniPtr.SetSimple()

	if game.FrameDSPtr != nil {
		game.FrameDSPtr.DeleteData() // free old GL buffers
	}

	game.FrameDSPtr = dsPtr

}

func (game *Mki3dGame) Update() {
	game.UpdateMonsters()
	game.UpdateTokens()
}

// Redraw the game stage
func (game *Mki3dGame) Redraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
	// draw stage
	game.StageDSPtr.SetBackgroundColor()
	game.StageDSPtr.DrawStage()
	// draw frame
	game.FrameDSPtr.DrawModel()
	// draw tokens
	// game.TokenDSPtr.DrawModel()
	game.DrawTokens()
	// draw monsters
	// game.MonsterDSPtr.DrawModel()
	game.DrawMonsters()

	// draw sectors
	gl.Disable(gl.DEPTH_TEST)
	game.SectorsDSPtr.DrawStage()
	gl.Enable(gl.DEPTH_TEST)

}

// RotHVType represents sequence of two rotations:
// by the angle XY on XY-plane and by the angle YZ on YZ-plane
// (in degrees)
type RotHVType struct {
	XZ float64
	YZ float64
}

type Traveler struct {
	Position mgl32.Vec3 // position
	Rot      RotHVType  // orientation
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
