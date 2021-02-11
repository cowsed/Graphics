package render

import (
	_ "fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	_ "image"
	_ "image/png"
	_ "math"
	_ "os"
	_ "time"
)

//Calculates the position in the 2d image world that is where the mouse is (2d because this does not yet get world position)
func CalculateGamePosition(win *pixelgl.Window, ScreenPos pixel.Vec) pixel.Vec {
	cam := pixel.IM.Scaled(camPos.Add(win.Bounds().Center()), camZoom).Moved(pixel.ZV.Sub(camPos))
	return cam.Unproject(ScreenPos).Scaled(2)
}

//Returns the first position chunkx,chunky,x,y,z that something is found
//Querying what is at that position with a tile vs actor is for later
func FindIntersect(basex, basey int) (int, int, int, int, int, bool) {
	//make list of possible places
	//Start pos = highest point
	for z := len(*(*ChunkReference)[0].WorldData) - heightCutoff - 1; z >= 0; z-- {
		x := basex + z
		y := basey - z


		if x >= 0 && y >= 0 {
			chunkx := x / 16
			chunky := y / 16
			ci := chunky*chunksDimension + chunkx
			//A thing is found
			if (*(*ChunkReference)[ci].WorldData)[z][y%16][x%16] != 0 {
				return chunkx, chunky, x, y, z, true
			}
		}
	}
	//error condition nothing found should have a better way of showing that though
	return 0, 0, 0, 0, 0, false
}

//Tells the renderer that something has changed
func SetChanged(change bool, index int) {
	(*ChunkReference)[index].SetChanged(change)
}

//Tells the renderer everything has changed
func SetAllChanged(val bool) {
	for i := 0; i < NumChunks; i++ {
		SetChanged(val, i)
	}
}

//Tell the renderer that the visuals must be updated
func SetDirty(change bool, index int) {
	(*ChunkReference)[index].SetDirty(change)
}

//Tells the renderer everything has changed
func SetAllDirty(val bool) {
	for i := 0; i < NumChunks; i++ {
		SetDirty(val, i)
	}
}

//Helper Loader Functions

//Bit of a debug function
//Draws grid lines at d
func DrawLines(w, h, d int, win *pixelgl.Window) {
	lineWidth := 2.0
	imd := imdraw.New(nil)

	for x := 0; x < w; x++ {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.RoundEndShape
		imd.Push(worldToIsoCoords(x-1, 0, d), worldToIsoCoords(x-1, h, d))
		imd.EndShape = imdraw.SharpEndShape
		imd.Line(lineWidth)
	}

	for y := 0; y < h; y++ {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.RoundEndShape
		imd.Push(worldToIsoCoords(-1, y, d), worldToIsoCoords(w-1, y, d))
		imd.EndShape = imdraw.SharpEndShape
		imd.Line(lineWidth)
	}
	imd.Draw(win)

}

func worldToIsoCoords(x, y, z int) pixel.Vec {
	//fmt.Println(x,y,z)
	return pixel.V(
		float64(((TileWidth/2)*x)+((TileHeight/2)*y)),
		float64((TileHeight/4)*(y-x)+z*(TileHeight/2)))
}

//Returns the bottom? level of what the world coords could be
//One would need to cast rays down the line to figure out which non-air object is being looked at
func isoToWorldCoords(v pixel.Vec) (int, int) {
	//Ah heck i  have to do a system of equation thing
	/*
		x'=(32.0*x)+(32*y)
		y'=(16.0*(y-x)))
		x'+2y'=64y
		y=(x'+2y')/64

		y'=16y-16x
		16x=16y-y'

	*/
	y := ((v.X + (2.0 * v.Y)) / 64)
	x := ((16.0*float64(y) - v.Y) / 16.0) + 1 //+1 to make things work or maybe not work but thats what it looks like
	//This messes up when it gets to the -x or -y gets off by 1
	if y < 0 {
		y--
	}
	if x < 0 {
		x--
	}

	return int(x), int(y)
}

func intMin(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
