package render

import (
	//"fmt"
	//"math/rand"
	"github.com/faiface/pixel"
	_ "github.com/faiface/pixel/imdraw"
	_ "github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	_ "golang.org/x/image/colornames"
	_ "image"
	_ "image/png"
)

//This pattern was stolen from a book
type ActorRenderer struct {
	Sheet          *pixel.Picture
	FrameIndex     int
	ChunkX, ChunkY int
	X, Y, Z        int
	Visible        bool
	
}

func (a ActorRenderer) makeKey() *[3]int {
	return &[3]int{a.X, a.Y, a.Z}
}
func (a ActorRenderer) makeChunkIndex() int {
	return a.ChunkX + a.ChunkY*chunksDimension
}

//Add Sprite adds the sprite to the pool of sprites to be included in the batch
//Key us an arguement because if the position changes the key changes so updating it doesnt work
func (a ActorRenderer) AddSprite() {
	println("Added a sprite")
	//Add the current sprite to the pool to be rendered
	//Blank key if calling externally

	//Create chunk index
	ci := a.makeChunkIndex()
	//Make a separate map for sprites to draw for each chunk.
	(*ChunkReference)[ci].AddSprite(&a)
	
	//This may be bad cuz it recalculates the max which could happen many times a frame (which is bad)
	(*ChunkReference)[ci].CalculateMax()

	SetChanged(true, ci)
}

//Remove the sprite from the pool
func (a ActorRenderer) RemoveSprite() {
	println("Removed a sprite")
	//Blank key if calling externally
	//Create chunk index
	ci := a.makeChunkIndex()

	(*ChunkReference)[ci].RemoveSprite(&a)


	SetChanged(true, ci)
}

func (a *ActorRenderer) UpdateVisibility(visible bool) {
	a.Visible = visible
	a.RemoveSprite()
	a.AddSprite()
}

func (a *ActorRenderer) UpdatePos(x, y, z int) {
	a.X = x
	a.Y = y
	a.Z = z
	a.RemoveSprite()
	a.AddSprite()
}
