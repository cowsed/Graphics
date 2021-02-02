package render

import (
	"fmt"
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	
	"image/color"

	_ "github.com/faiface/pixel/text"
	_ "image"
	_ "image/png"
	_ "os"
)

//VoidColor defines the color of the background
var VoidColor color.RGBA = colornames.Skyblue




var spriteSheet pixel.Picture //maybe load this by stitching together others but for now this
var sheetFrames []pixel.Rect

var TileBatches []*pixel.Batch //Set of Sprites

var SpritesToDraw []map[string]*ActorRenderer

var changes = []bool{true,true,true,true,true,true,true,true,true}
//var mapUpdated = true



//Render Everything
func Render(win *pixelgl.Window, WorldMap []*[][][]int, w,h int){
	drawStart:=time.Now()

	//Calculate camera positioning and UI positioning
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(pixel.ZV.Sub(camPos))
	//oppCam := pixel.IM.Moved(camPos).Scaled(camPos, 1/camZoom)
	win.SetMatrix(cam)

	//Clear the Window to prepare for drawing
	win.Clear(VoidColor)
	RenderWorld(win, WorldMap ,w,h)
	//render.DrawLines(len(WorldMap[0][0]),len(WorldMap[0]),4,win)
	
	//Render UI
	RenderUI(win)

	//Timing things
	drawDt := time.Since(drawStart).Seconds()
	SendString(fmt.Sprintf("Render Time(ms): %f\n",1000*drawDt/60.0))
}

func RenderWorld(win *pixelgl.Window, Chunks []*[][][]int , w,h int){
	//len of chunks should be 9. 3x3 grid
	x:=0
	y:=0
	for c:=0; c<9; c++{
		y=3-c/3
		x=c%3
		
		RenderChunk(win, Chunks, &changes[c], w,h,x,y, c)
	}
}

func RenderChunk(win *pixelgl.Window, WorldMaps []*[][][]int , changed *bool, w,h,chunkX,chunkY, ChunkIndex int) {
	//Trackers
	spritesDrawn := 0
	tilesDrawn := 0
	//Create Chunk Map
	ChunkMap:=WorldMaps[ChunkIndex]
	//Dimensions of the world map
	d := len((*ChunkMap))
	//h := len((*WorldMap)[0])
	//w := len((*WorldMap)[0][0])

	if *changed { //Reset Batch

		TileBatches[ChunkIndex].Clear()

		for z := 0; z < d-heightCutoff; z++ {
			for y := h - 1; y >= 0; y-- { //Go backwards to go over each other
				yp := h - y - 1 //needs to be reversed because of orientation
				for x := 0; x < w; x++ {

					if visible, inside := CheckVisibility(x, y, yp, z, w, h, d, heightCutoff, ChunkIndex, WorldMaps); visible {
						//if its inside darken it ab it
						if inside{
							TileBatches[ChunkIndex].SetColorMask(color.RGBA{60,60,60,255})
						}
						//Render Sprites (This is sort of a bad idea because it takes a map which is unfun to allocate)
						//But may be better than searching through the list of
						key := string(x) + "" + string(y) + "," + string(z)
						if sprite, ok := SpritesToDraw[ChunkIndex][key]; ok { //Check if there is a Sprite here
							//Render the sprite
							var spriteIndex int

							if (*sprite).Visible{
								spriteIndex=(*sprite).FrameIndex
							} else {
								spriteIndex=159
							}
								mx := pixel.IM.Moved(toIsoCoords(x+(chunkX*w), y+(chunkY*h), z)) //Position sprite in space
								sprite := pixel.NewSprite(spriteSheet, sheetFrames[spriteIndex])
								sprite.Draw(TileBatches[ChunkIndex], mx)
								spritesDrawn++

						}
						//fmt.Println(x+(chunkX*w), y+(chunkY*h))
						//Render World
						//Position
						mx := pixel.IM.Moved(toIsoCoords(x+(chunkX*w), y+(chunkY*h), z))
						//Material Index
						TileIndex := (*ChunkMap)[z][yp][x]
						if TileIndex != 0 {
							TileIndex-- //Decrement to index

							
							sprite := pixel.NewSprite(spriteSheet, sheetFrames[TileIndex])
							sprite.Draw(TileBatches[ChunkIndex], mx)
							tilesDrawn++
						}
						if inside{
							TileBatches[ChunkIndex].SetColorMask(color.RGBA{255,255,255,255})
						}
					}

				}
			}
		}

	}

	TileBatches[ChunkIndex].Draw(win) //Draw the batch no matter what

	*changed=false //Reset changes



}

func CheckVisibility(x, y, yp, z, w, h, d, z_cutoff, ChunkIndex int, ChunkMaps []*[][][]int) (bool,bool) {
	onFrontFace := (x == w-1) || (y == 0) || (z == d-z_cutoff-1)
	
	onFullXFace:=(ChunkIndex/3==2) && y==0
	onFullYFace:=(ChunkIndex%3==2) && x==w-1
	onFullXYFace:=onFullXFace||onFullYFace
	
	onXYFace := ((x == w-1) || (y == 0))
	exposed := false
	if z != d-1 {
		yp := h - y - 1
		exposed = (*ChunkMaps[ChunkIndex])[z+1][yp][x] == 0
	}
	inside:=!exposed
	var exposedToOtherChunkAir bool
	if inside && onXYFace{ //If it is on a side visible to the player and is considered inside based off of first look up
		//Check if its actually inside based in look into other chunks
		//Check further y forward
		if ChunkIndex/3 <2{ //It is not the last one
			//Check the next chunk at the first one at the same x and z
			exposedToOtherChunkAir=(*ChunkMaps[ChunkIndex+3])[z][0][x] == 0
		}
		//Check x further right
		if ChunkIndex%3 <2{ //It is not the furthest right
			//Check the next chunk ath the first one at the same x and y
			exposedToOtherChunkAir = exposedToOtherChunkAir || (*ChunkMaps[ChunkIndex+1])[z][yp][0] == 0
		}
		
	}
	
	return onFrontFace || exposed,  !onFullXYFace && !(exposed||exposedToOtherChunkAir) 
}

func InitRender() {
	spriteSheet, sheetFrames = loadSheet("Assets/outside4.png", 64, 64)
	fmt.Println("Loaded Environment")
	
	//Create the batches and the sprite maps
	SpritesToDraw = make([]map[string]*ActorRenderer, 9)
	TileBatches=make([]*pixel.Batch,9)
	for i:=0; i<9; i++{
		SpritesToDraw[i] = make(map[string]*ActorRenderer)
		TileBatches[i] = pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)
	}
}
