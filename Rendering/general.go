package render

import (
	_"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
)
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


func DrawLines(w, h, d int, win *pixelgl.Window) {
	lineWidth := 2.0
	imd := imdraw.New(nil)

	for x := 0; x < w; x++ {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.RoundEndShape
		imd.Push(toIsoCoords(x-1, 0, d), toIsoCoords(x-1, h, d))
		imd.EndShape = imdraw.SharpEndShape
		imd.Line(lineWidth)
	}

	for y := 0; y < h; y++ {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.RoundEndShape
		imd.Push(toIsoCoords(-1, y, d), toIsoCoords(w-1, y, d))
		imd.EndShape = imdraw.SharpEndShape
		imd.Line(lineWidth)
	}
	imd.Draw(win)

}

func toIsoCoords(x, y, z int) pixel.Vec {
	return pixel.V(
		float64(32.0*(x))+float64(32*y),
		float64(16.0*(y-x)+z*32))
}
