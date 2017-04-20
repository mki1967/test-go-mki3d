package main

import (
	"fmt" // tests
	// "errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	// "github.com/go-gl/glfw/v3.2/glfw"
	// "math/rand"
)

const BoxMargin = 40 // margin for bounding box of the stage

var FrameColor = mki3d.Vector3dType{1.0, 1.0, 1.0} // color of the bounding box frame

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

	return nil
}

// Load token shape and init the tokenDSPtr.
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

	return nil
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

// Redraw the game stage
func (game *Mki3dGame) Redraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
	// draw stage
	game.StageDSPtr.SetBackgroundColor()
	game.StageDSPtr.DrawStage()
	// draw frame
	game.FrameDSPtr.DrawModel()
	// draw tokens
	game.TokenDSPtr.DrawModel()
	// draw monsters
	game.MonsterDSPtr.DrawModel()

	// draw sectors
	gl.Disable(gl.DEPTH_TEST)
	game.SectorsDSPtr.DrawStage()
	gl.Enable(gl.DEPTH_TEST)

}

type Traveler struct {
	Position mgl32.Vec3 // position
	/* orientation */
	rotXZ float32 // horizontal rotation (in degrees)
	rotYZ float32 // vertical rotation (in degrees)

}
