package render

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"

	"image/color"

	_ "github.com/faiface/pixel/text"
	_ "image"
	_ "image/png"
	_ "os"
)

//NumChunks defines the number of chunks visible at once
const NumChunks = 25
const chunksDimension = 5

//VoidColor defines the color of the background
var VoidColor color.RGBA = colornames.Skyblue

var spriteSheet pixel.Picture //maybe load this by stitching together others but for now this
var sheetFrames []pixel.Rect

var TileBatches []*pixel.Batch //Set of Sprites

var SpritesToDraw []map[string]*ActorRenderer

//Sprite used for the selection cursor
var SelectSprite *pixel.Sprite

var ChunkReference *[]Chunk

//Render Everything
func Render(win *pixelgl.Window, w, h,d int) {

	//Show Height Cutoff
	SendString(fmt.Sprintf("Height Cutoff: %d\n", d-heightCutoff))

	//Calculate camera positioning and UI positioning
	cam := pixel.IM.Scaled(camPos.Add(win.Bounds().Center()), camZoom).Moved(pixel.ZV.Sub(camPos))
	//oppCam := pixel.IM.Moved(camPos).Scaled(camPos, 1/camZoom)
	win.SetMatrix(cam)

	//Clear the Window to prepare for drawing
	win.Clear(VoidColor)
	worldStart:=time.Now()
	RenderWorld(win, w, h,d)
	SendString(fmt.Sprintf("World Render Time(ms): %d\n",time.Since(worldStart).Milliseconds()))


	//Render Selectioncursor DB
	//cursorStart:=time.Now()
	//x, y := isoToWorldCoords(CalculateGamePosition(win, win.MousePosition()))
	//mx := pixel.IM.Moved(worldToIsoCoords(x, y, 0))
	//SelectSprite.Draw(win, mx)
	//SendString(fmt.Sprintf("Cursor Time(ms): %d\n",time.Since(cursorStart).Milliseconds()))


	//lineStart:=time.Now()
	//DrawLines(32, 32, 0, win)
	//SendString(fmt.Sprintf("Line Draw Time(ms): %d\n",time.Since(lineStart).Milliseconds()))

	//Render UI
	RenderUI(win)

}

func RenderWorld(win *pixelgl.Window,  w, h,d int) {
	//last := time.Now()
	//len of chunks should be 9. 3x3 grid
	x := 0
	y := 0
	for c := 0; c < NumChunks; c++ {
		y = chunksDimension - c/chunksDimension
		x = c % chunksDimension

		RenderChunk(win, w, h,d, x, y, c)
	}

}

func RenderChunk(win *pixelgl.Window, w, h,d, chunkX, chunkY, ChunkIndex int) {
	//Trackers
	spritesDrawn := 0
	tilesDrawn := 0
	//Create Chunk Map
	ChunkMap := (*ChunkReference)[ChunkIndex]
	//Dimensions of the world map

	if ChunkMap.GetChanged() { //Reset Batch

		TileBatches[ChunkIndex].Clear()

		for z := 0; z < ChunkMap.MaxHeight-heightCutoff; z++ {
			for y := h - 1; y >= 0; y-- { //Go backwards to go over each other
				//yp := h - y - 1 //needs to be reversed because of orientation
				for x := 0; x < w; x++ {

					if visible, inside := CheckVisibility(x, y, z, w, h, d, heightCutoff, ChunkIndex); visible {
						//if its inside darken it ab it
						if inside {
							TileBatches[ChunkIndex].SetColorMask(color.RGBA{60, 60, 60, 255})
						}
						//Render Sprites (This is sort of a bad idea because it takes a map which is unfun to allocate)
						//But may be better than searching through the list of
						key := string(x) + "" + string(y) + "," + string(z)
						if sprite, ok := SpritesToDraw[ChunkIndex][key]; ok { //Check if there is a Sprite here
							//Render the sprite
							var spriteIndex int

							if (*sprite).Visible {
								spriteIndex = (*sprite).FrameIndex
							} else {
								spriteIndex = 159
							}
							mx := pixel.IM.Moved(worldToIsoCoords(x+(chunkX*w), y+(chunkY*h), z)) //Position sprite in space
							sprite := pixel.NewSprite(spriteSheet, sheetFrames[spriteIndex])
							sprite.Draw(TileBatches[ChunkIndex], mx)
							spritesDrawn++

						}
						//fmt.Println(x+(chunkX*w), y+(chunkY*h))
						//Render World
						//Position
						mx := pixel.IM.Moved(worldToIsoCoords(x+(chunkX*w), y+(chunkY*h), z))
						//Material Index
						TileIndex := (*ChunkMap.WorldData)[z][y][x]
						if TileIndex != 0 {
							TileIndex-- //Decrement to index

							sprite := pixel.NewSprite(spriteSheet, sheetFrames[TileIndex])
							sprite.Draw(TileBatches[ChunkIndex], mx)
							tilesDrawn++
						}
						if inside {
							TileBatches[ChunkIndex].SetColorMask(color.RGBA{255, 255, 255, 255})
						}
					}

				}
			}
		}

	}

	TileBatches[ChunkIndex].Draw(win) //Draw the batch no matter what

	(*ChunkReference)[ChunkIndex].SetChanged(false) //Reset changes

}

func CheckVisibility(x, y, z, w, h, d, z_cutoff, ChunkIndex int) (bool, bool) {
	onFrontFace := (x == w-1) || (y == 0) || (z == d-z_cutoff-1)

	onFullXFace := (ChunkIndex/chunksDimension == chunksDimension-1) && y == 0
	onFullYFace := (ChunkIndex%chunksDimension == chunksDimension-1) && x == w-1
	onFullXYFace := onFullXFace || onFullYFace

	onXYFace := ((x == w-1) || (y == 0))
	exposed := false
	if z != d-1 {
		exposed = (*(*ChunkReference)[ChunkIndex].WorldData)[z+1][y][x] == 0
	}
	inside := !exposed
	var exposedToOtherChunkAirY bool
	var exposedToOtherChunkAirX bool

	if inside && onXYFace { //If it is on a side visible to the player and is considered inside based off of first look up
		//Check if its actually inside based in look into other chunks
		//Check further y forward
		if ChunkIndex/chunksDimension < chunksDimension-1 { //It is not the last one
			//Check the next chunk at the first one at the same x and z
			exposedToOtherChunkAirY = (*(*ChunkReference)[ChunkIndex+chunksDimension].WorldData)[z][0][x] == 0
		}
		//Check x further right
		if ChunkIndex%chunksDimension < chunksDimension-1 { //It is not the furthest right
			//	//Check the next chunk ath the first one at the same y and z
			exposedToOtherChunkAirX = (*(*ChunkReference)[ChunkIndex+1].WorldData)[z][y][0] == 0
			//	exposedToOtherChunkAir=true
		}

	}

	return onFrontFace || exposed, !onFullXYFace && !(exposed || exposedToOtherChunkAirY || exposedToOtherChunkAirX)
}

func InitRender() {
	spriteSheet, sheetFrames = loadSheet("Assets/outside4.png", 64, 64)
	fmt.Println("Loaded Environment")

	//Create the batches and the sprite maps
	SpritesToDraw = make([]map[string]*ActorRenderer, NumChunks)
	TileBatches = make([]*pixel.Batch, NumChunks)
	for i := 0; i < NumChunks; i++ {
		SpritesToDraw[i] = make(map[string]*ActorRenderer)
		TileBatches[i] = pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)
	}
	
	SelectSprite = pixel.NewSprite(spriteSheet, sheetFrames[159])
}
