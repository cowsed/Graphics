package render

import (
	"fmt"
	//"math/rand"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"image"
	_ "image/png"
	"os"
)

//Global variables that will be referenced
var sheet pixel.Picture
var frames []pixel.Rect

var EnvBatch *pixel.Batch

//Loads tiles
func BackgroundInit() {

	sheet, frames = loadSheet("Assets/outside4.png",64,64)
	fmt.Println("Loaded Environment")
}

func DrawMap(win *pixelgl.Window, EnvMap [][][]int, changed bool, heightCutoff int,dbText *text.Text) {

	d := len(EnvMap)
	h := len(EnvMap[0])
	w := len(EnvMap[0][1])


	spritesDrawn := 0

	//Optimization usung batch - get a better tileset so this weirdness inst necessary


	if changed{
		 EnvBatch=pixel.NewBatch(&pixel.TrianglesData{}, sheet)

		EnvBatch.Clear()

		for z := 0; z < d-heightCutoff; z++ {
			for y := h - 1; y >= 0; y-- { //Go backwards to go over each other
				for x := 0; x < w; x++ {
					yp := h - y - 1
					mx := pixel.IM.Moved(toIsoCoords(x,y,z))
	
					TileIndex := EnvMap[z][yp][x]
					if TileIndex > 0 { //Non Void tile
						TileIndex -= 1 //drop down to actual index cuz void is 0
	
						//Check if it has a tile blocking it
						exposed := true
						onFrontFace := (x == w-1) || (y == 0) || (z == d-heightCutoff-1)
	
						
						//Check for others
						//onLastFace:=(x==0)||(y==0)||(z==0)
						//if !onLastFace && !onFrontFace{
						//	exposed=(EnvMap[x][y][z+1]==0)
						//}
						if z != d-1 {
							exposed = EnvMap[z+1][yp][x] == 0
						}
						if onFrontFace || exposed { //If Tile is too be drawn
							sprite := pixel.NewSprite(sheet, frames[TileIndex])
							sprite.Draw(EnvBatch, mx)
	
							spritesDrawn++
						}
					}
				}
			}
		}
				fmt.Println( "Sprites: ", spritesDrawn)
	}
		EnvBatch.Draw(win)
	
		fmt.Fprintln(dbText, "Sprites: ", spritesDrawn)


}

func loadSheet(fname string, w,h int) (pixel.Picture, []pixel.Rect) {

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

/*


const (
	NONE = iota
	DIRT_GRASS
	DIRT_GRASS_FRONT_RAMP
	DIRT_GRASS_SIDE_RAMP
	DIRT_CHECKERS

	GRASS_GRASS
	GRASS_GRASS_FRONT_RAMP
	GRASS_GRASS_SIDE_RAMP
	GRASS_GRASS_CHECKERS

	GRASS_SNOW
	GRASS_SNOW_FRONT_RAMP
	GRASS_SNOW_SIDE_RAMP

	SNOW_SNOW
	SNOW_SNOW_FRONT_RAMP
	SNOW_SNOW_SIDE_RAMP

	SNOW_GRASS
	SNOW_GRASS_CHECKERS
	SNOW_GRASS_FRONT_RAMP
	SNOW_GRASS_SIDE_RAMP

	DARK_LIGHT_STONE
	DARK_LIGHT_STONE_FRONT_RAMP
	DARK_LIGHT_STONE_SIDE_RAMP

	DARK_DARK_STONE_BLOCK
	DARK_DARK_STONE
	DARK_DARK_STONE_FRONT_RAMP
	DARK_DARK_STONE_SIDE_RAMP

)


func makeFnames(n int) []string {
	Ss := make([]string, 0)
	for i := 1; i <= n; i++ {
		Ss = append(Ss, fmt.Sprintf("Assets/Blocks/blocks_%d.png", i))
	}
	return Ss
}

func loadTileSet() []pixel.Picture {
	m := 101
	out := make([]pixel.Picture, m)

	Ss := makeFnames(m)
	for i, v := range Ss {
		fmt.Println("Pic: ", i, " - ", v)
	}
	for i, s := range Ss {
		pic, err := loadPicture(s)
		if err != nil {
			panic(err)
		}
		out[i] = pic

	}
	return out
}

*/

//the
/*
	EnvMap := [2][9][12]int{
	[9][12]int{

			[12]int{1, 1, 1, 1, 1, 82, 1, 82,88,1,1,52},
			[12]int{88, 1, 1, 1, 1, 1, 1, 1,1,1,1,52},
			[12]int{1, 1, 79, 1, 1, 88, 1, 1,1,91,1,52},
			[12]int{1, 1, 1, 1, 1, 1, 1, 91,1,1,1,52},
			[12]int{1, 12, 92, 1, 90, 76, 1, 1,1,1,1,52},
			[12]int{1, 1, 1, 1, 1, 1, 1, 1,91,1,1,52},
			[12]int{1, 1, 1, 90, 1, 1, 89, 1,1,78,1,52},
			[12]int{1, 73, 1, 1, 1, 1, 1, 1,1,77,1,52},
			[12]int{52, 0, 0, 1, 96, 96, 96, 96,96,96,96,52}, },


		[9][12]int{
			[12]int{63,0,0,0,0,0,0,1,3,0,0,0},
			[12]int{63,0,0,0,0,0,0,2,0,0,0,0},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,1},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,1},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,1},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,0},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,0},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,0},
			[12]int{63,0,0,0,0,0,0,0,0,0,0,0},
		},
	}
*/
