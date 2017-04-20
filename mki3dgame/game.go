package main

import (
	// "fmt" // tests
	// "errors"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/test-go-mki3d/tmki3d"
	// "github.com/go-gl/glfw/v3.2/glfw"
	// "math/rand"
)

const MARGIN = 40 // margin for bounding box of the stage

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

	VMin, VMax mgl32.Vec3 // corners of the bounding box of the stage (computed with the MARGIN)

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

	game.TokenDSPtr = tokenDataShaderPtr

	return nil
}

// Load sectors shape and init the related data.
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

	game.StageDSPtr = stageDataShaderPtr

	return nil
}

// Redraw the game stage
func (game *Mki3dGame) Redraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
	// draw stage
	game.StageDSPtr.SetBackgroundColor()
	game.StageDSPtr.DrawStage()
	// draw tokens
	game.TokenDSPtr.DrawModel()

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
