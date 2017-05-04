This is a game that demonstrates the use of the data produced with MKI3D web 3D editor (see: https://mki1967.github.io/mki3d/ ).
In the game, you have to collect tokens scattered in the stages and avoid being captured by the monsters.
(A short screen-cast is available at: https://youtu.be/vp6nhvOqhdU . )
Run the game with the path to assets directory as the command line argument.
(See the content of the runme script in this directory.)
The assets directory has the following general structure:
assets
├── icons
│   └── ... (icon '.png' files)
├── monsters
│   └── ... (monster shapes '.mki3d' files - made with MKI3D)
├── sectors
│   └── ... (shapes of screen sectors '.mki3d' - made with MKI3D, specific to the code )  
├── stages
│   └── ... (stages '.mki3d' files - made with MKI3D)
└── tokens
    └── ... (token shapes '.mki3d' files - made with MKI3D)

You can design your own stages and the shapes of monsters or tokens
with this editor.
Just place the files in the respective sub-directories
'stages', 'monsters', or 'tokens' of the main assets directory.
Shapes are selected randomly from each sub-directory for each stage.
At this moment the included 'assets' directory has the following contents: 

assets
├── icons
│   ├── mkisg_icon_16x16.png
│   ├── mkisg_icon_32x32.png
│   └── mkisg_icon_48x48.png
├── monsters
│   └── monster_1.mki3d
├── sectors
│   └── sectors.mki3d
├── stages
│   ├── stage10.mki3d
│   ├── stage1.mki3d
│   ├── stage-2016-08-10.mki3d
│   ├── stage-2016-11-14.mki3d
│   ├── stage-2017-04-02.mki3d
│   ├── stage2.mki3d
│   ├── stage3.mki3d
│   ├── stage4.mki3d
│   ├── stage5.mki3d
│   ├── stage6.mki3d
│   ├── stage7.mki3d
│   ├── stage8.mki3d
│   └── stage9.mki3d
└── tokens
    └── item.mki3d
