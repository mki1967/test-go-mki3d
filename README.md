# test-go-mki3d

This repository contains programs for demonstrating
go-mki3d packages from the repository:
https://github.com/mki1967/go-mki3d.git
that contain tools for reading and displaying the shapes designed with
MKI3D editor (see: https://mki1967.github.io/mki3d/ ) in Go language with the use of
"https://github.com/go-gl" libraries.

At this moment the repository contains:

* simple - a simple demo of reading and displaying  MKI3D data.
* mki3dgame - a game, where you have to collect tokens on randomly selected stages and avoid being captured by the monsters
  (the stages, tokens and monsters are designed with MKI3D and just placed in specific sub-directories. You can design your own stages.)

To run the programs you need the Go language  and the following packages:

*	"github.com/go-gl/gl/v3.3-core/gl"
*	"github.com/go-gl/glfw/v3.2/glfw"
*	"github.com/go-gl/mathgl/mgl32"
*	"github.com/mki1967/go-mki3d/mki3d"
*	"github.com/mki1967/go-mki3d/glmki3d"





