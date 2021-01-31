package render

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	_ "image"
	_ "image/png"
	_ "os"
)

var spriteSheet pixel.Picture //maybe load this by stitching together others but for now this
var sheetFrames []pixel.Rect

var TileBatch *pixel.Batch //Set of Sprites

var SpritesToDraw map[ string ]*ActorRenderer



func RenderAll(win *pixelgl.Window, WorldMap *[][][]int, changed bool, heightCutoff int) {
	//Trackers
	spritesDrawn := 0
	tilesDrawn:=0
	//Dimensions of the world map
	d := len((*WorldMap))
	h := len((*WorldMap)[0])
	w := len((*WorldMap)[0][1])

	if changed { //Reset Batch
		TileBatch = pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)
		TileBatch.Clear()
		
		for z := 0; z < d-heightCutoff; z++ {
			for y := h - 1; y >= 0; y-- { //Go backwards to go over each other
				yp := h - y - 1 //needs to be reversed because of orientation
				for x := 0; x < w; x++ {
				
					if CheckVisibility(x, y, yp, z, w, h, d, heightCutoff, WorldMap) {
						//Render Sprites (This is sort of a bad idea because it takes a map which is unfun to allocate)
						//But may be better than searching through the list of 
						key:=string(x)+""+string(y)+","+string(z)
						if val, ok := SpritesToDraw[key]; ok {//Check if there is a Sprite here
    						//Render the sprite
							mx := pixel.IM.Moved(toIsoCoords(x, y, z)) //Position sprite in space
							sprite := pixel.NewSprite(spriteSheet, sheetFrames[(*val).FrameIndex])
							sprite.Draw(TileBatch, mx)
							spritesDrawn++
						}
						
						//Render World
						//Position
						mx := pixel.IM.Moved(toIsoCoords(x, y, z))
						//Material Index
						TileIndex := (*WorldMap)[z][yp][x]
						if TileIndex!=0{
							TileIndex-- //Decrement to index
							sprite := pixel.NewSprite(spriteSheet, sheetFrames[TileIndex])
							sprite.Draw(TileBatch, mx)
							tilesDrawn++
						}
					}

				}
			}
		}

	}
	TileBatch.Draw(win) //Draw the batch no matter what
	
}

func CheckVisibility(x, y, yp, z, w, h, d, z_cutoff int, WorldMap *[][][]int) bool {
	onFrontFace := (x == w-1) || (y == 0) || (z == d-z_cutoff-1)

	exposed := false
	if z != d-1 {
		yp := h - y - 1
		exposed = (*WorldMap)[z+1][yp][x] == 0
	}

	return onFrontFace || exposed
}

func InitRender() {
	spriteSheet, sheetFrames = loadSheet("Assets/outside4.png", 64, 64)
	fmt.Println("Loaded Environment")
	SpritesToDraw=make( map[ string ]*ActorRenderer)

}
