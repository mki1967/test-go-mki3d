package main

import (
	"io/ioutil"
	// "fmt"
	// "errors"
	"github.com/mki1967/go-mki3d/mki3d"
	"math/rand"
	"os"
)

// the structure for assets of mki3dgame
type Assets struct {
	Path     string
	Assets   []os.FileInfo
	Stages   []os.FileInfo
	Tokens   []os.FileInfo
	Monsters []os.FileInfo
	Icons    []os.FileInfo
}

const (
	StagesDir = "stages"
)

func LoadAssets(pathToAssets string) (*Assets, error) {
	ass, err := ioutil.ReadDir(pathToAssets)
	if err != nil {
		return nil, err
	}

	assets := Assets{Path: pathToAssets, Assets: ass} /// ...

	assets.Stages, err = ioutil.ReadDir(pathToAssets + "/" + StagesDir)

	if err != nil {
		return &assets, err
	}

	return &assets, nil
}

func (a *Assets) LoadRandomStage() (*mki3d.Mki3dType, error) {

	r := rand.Intn(len(a.Stages))

	mki3dPtr, err := mki3d.ReadFile(a.Path + "/" + StagesDir + "/" + a.Stages[r].Name())
	if err != nil {
		return nil, err
	}

	return mki3dPtr, nil

}
