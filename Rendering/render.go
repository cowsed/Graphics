package render

import (
	"../Materials"
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

//NumChunks defines the number of chunks visible at once
const NumChunks = 25
const chunksDimension = 5

//VoidColor defines the color of the background
var VoidColor color.RGBA = colornames.Skyblue

//Stuff for sprites
//tile dimensions. should probably always be the same
const TileWidth = 32
const TileHeight = 32

var spriteSheet pixel.Picture //maybe load this by stitching together others but for now this

var sheetFrames []pixel.Rect

//Sprite used for the selection cursor
var SelectSprite *pixel.Sprite

var ChunkReference *[]Chunk

//Render Everything
func Render(win *pixelgl.Window, w, h, d int) {

	//Show Height Cutoff
	SendString(fmt.Sprintf("Height Cutoff: %d\n", heightCutoff))

	//Calculate camera positioning and UI positioning
	cam := pixel.IM.Scaled(camPos.Add(win.Bounds().Center()), camZoom).Moved(pixel.ZV.Sub(camPos))
	win.SetMatrix(cam)

	//Clear the Window to prepare for drawing
	win.Clear(VoidColor)
	worldStart := time.Now()
	RenderWorld(win, w, h, d)
	SendString(fmt.Sprintf("World Render Time(ms): %d\n", time.Since(worldStart).Milliseconds()))

	//Render Selectioncursor DB
	cursorStart := time.Now()
	x, y := isoToWorldCoords(CalculateGamePosition(win, win.MousePosition()))

	chunkx, chunky, x2, y2, z2, success := FindIntersect(x, y)

	if success {
		seenTileIndex := (*(*ChunkReference)[chunky*5+chunkx].WorldData)[z2][y2%16][x2%16] - 1
		SendString(fmt.Sprintf("TileIndex %d : %s\n", seenTileIndex, materials.SpritesByIndex[seenTileIndex]))
		SendString(fmt.Sprintln("Desc: ",materials.Descriptions[seenTileIndex]))

	} else {
		SendString(fmt.Sprintln("TileIndex X : No Block found"))
	}

	mx := pixel.IM.Moved(worldToIsoCoords(x2, y2, z2))
	SelectSprite.Draw(win, mx)

	mx2 := pixel.IM.Moved(worldToIsoCoords(x, y, 0))
	SelectSprite.Draw(win, mx2)

	SendString(fmt.Sprintf("Cursor Time(ms): %.3f\n", time.Since(cursorStart).Seconds()*1000))

	//lineStart:=time.Now()
	//DrawLines(32, 32, -1, win)
	//SendString(fmt.Sprintf("Line Draw Time(ms): %d\n",time.Since(lineStart).Milliseconds()))

	//Render UI
	RenderUI(win)

}

func RenderWorld(win *pixelgl.Window, w, h, d int) {
	x := 0
	y := 0
	for c := NumChunks - 1; c >= 0; c-- {
		y = c / chunksDimension
		x = c % chunksDimension

		ci := y*chunksDimension + (chunksDimension - x - 1)
		//This will only update if it's been marked necessary - cool huh
		(*ChunkReference)[ci].UpdateTiles(c, heightCutoff)
		(*ChunkReference)[ci].Render(win, chunksDimension-x, y)
	}

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
	//TODO FIX THis so it works with sensible chunk coordinates
	return true || onFrontFace || exposed, !(true || !onFullXYFace && !(exposed || exposedToOtherChunkAirY || exposedToOtherChunkAirX))
}

//Initialize the render
func InitRender() {
	//spriteSheet, sheetFrames = loadSheet("Assets/Custom3.png", TileWidth, TileHeight)
	spriteSheet, sheetFrames = materials.GetData()

	fmt.Println("Loaded Environment")

	//Create the batches and the sprite maps

	for i := 0; i < NumChunks; i++ {
		(*ChunkReference)[i].Init()
	}

	SelectSprite = pixel.NewSprite(spriteSheet, sheetFrames[1])
}
