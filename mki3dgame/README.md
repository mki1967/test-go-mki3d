NOTE: The game mki3dgame will be further developed as a standalone project at https://github.com/mki1967/mki3dgame


This is a game that demonstrates the use of the data produced with MKI3D web 3D editor (see: https://mki1967.github.io/mki3d/ ).

To run the game you need the packages:
*	"github.com/go-gl/gl/v3.3-core/gl"
*	"github.com/go-gl/glfw/v3.2/glfw"
*	"github.com/go-gl/mathgl/mgl32"
*	"github.com/mki1967/go-mki3d/mki3d"
*	"github.com/mki1967/go-mki3d/glmki3d"

In the game, you have to collect tokens scattered in the stages and avoid being captured by the monsters.
(A short screen-cast is available at: https://youtu.be/vp6nhvOqhdU . )
Run the game with the path to assets directory as the command line argument.
(See the content of the runme script in this directory.)
The assets directory has the following sub-directories:

* 'icons' -  icon '.png' files (some systems may use them ...)
* 'monsters' - monster shapes '.mki3d' files - made with MKI3D
* 'sectors'  - shapes of screen sectors '.mki3d' - made with MKI3D, specific to the code 
* 'stages'  - stages '.mki3d' files - made with MKI3D
* 'tokens'  - token shapes '.mki3d' files - made with MKI3D

You can design your own stages and the shapes of monsters or tokens
with this editor.
Just place the files in the respective sub-directories
'stages', 'monsters', or 'tokens' of the main assets directory.
Shapes are selected randomly from each sub-directory for each stage.
