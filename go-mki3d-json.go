package main

import (
	"fmt"
	"github.com/mki1967/go-mki3d/mki3d"
	"encoding/json"
)


func main() {


	mki3dData, err := mki3d.ReadFile("noname.mki3d")
	if err!=nil {
		panic( err )
	}
	fmt.Println(mki3dData)
	mki3dC, err := json.Marshal(mki3dData)
	if err!=nil {
		panic( err )
	}
	fmt.Println(string(mki3dC))

}
