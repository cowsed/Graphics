package render

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	_ "math"
	"os"
	"time"
)

//Calculates the position in the 2d image world that is where the mouse is (2d because this does not yet get world position)
func CalculateGamePosition(win *pixelgl.Window, ScreenPos pixel.Vec) pixel.Vec {
	cam := pixel.IM.Scaled(camPos.Add(win.Bounds().Center()), camZoom).Moved(pixel.ZV.Sub(camPos))
	return cam.Unproject(ScreenPos)
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

//Loads a sprite sheet into chunks of wxh
//Returns the loaded image and the bounds of the sprites
func loadSheet(fname string, w, h int) (pixel.Picture, []pixel.Rect) {

	spritesheet, err := loadPicture(fname)
	if err != nil {
		panic(err)
	}

	var frames []pixel.Rect
	for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += float64(h) {
		for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += float64(w) {
			frames = append(frames, pixel.R(x, y, x+float64(w), y+float64(h)))
		}
	}

	return spritesheet, frames
}

//Helper Loader Functions

//Loads the picture
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

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
	return pixel.V(
		float64(32.0*(x))+float64(32*y),
		float64(16.0*(y-x)+z*32))
}

//Returns the bottom? level of what the world coords could be
//One would need to cast rays down the line to figure out which non-air object is being looked at
//TODO Fix this because it does not track out correctly
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
	isoConversionTime := time.Now()
	y := ((v.X + (2.0 * v.Y)) / 64)
	x := ((16.0*float64(y) - v.Y) / 16.0) + 1 //+1 to make things work or maybe not work but thats what it looks like
	//This messes up when it gets to the -x or -y gets off by 1
	if y < 0 {
		y--
	}
	if x < 0 {
		x--
	}

	SendString(fmt.Sprintf("Iso Conversion Time(ms): %f\n", time.Since(isoConversionTime).Seconds()*1000))
	return int(x), int(y)
}

func intMin(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
