package main

import (
	"errors"
	"fmt"
	"github.com/mki1967/go-mki3d/mki3d"
	"os"
)

func main() {

	// get file name from command line argument
	if len(os.Args) < 2 {
		panic(errors.New(" *** PROVIDE FILE NAME AS A COMMAND LINE ARGUMENT !!! *** "))
	}

	// Load mki3d data from a file
	mki3dPtr, err := mki3d.ReadFile(os.Args[1]) // load mki3d data from the file produced by MKI3D
	if err != nil {
		panic(err)
	}

	fmt.Println(mki3d.Stringify(mki3dPtr))

}
