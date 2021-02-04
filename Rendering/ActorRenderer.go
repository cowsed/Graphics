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
func (a ActorRenderer) AddSprite(key *[3]int) {
	//Add the current sprite to the pool to be rendered
	//Blank key if calling externally
	if key == nil {
		key = a.makeKey()
	}
	//Create chunk index
	ci := a.makeChunkIndex()
	//Make a separate map for sprites to draw for each chunk.
	(*(*ChunkReference)[ci].SpriteData)[*key] = &a
	(*ChunkReference)[ci].CalculateMax()

	SetChanged(true, ci)
}

//Remove the sprite from the pool
func (a ActorRenderer) RemoveSprite(key *[3]int) {
	//Blank key if calling externally
	if key == nil {
		key = a.makeKey()
	}
	//Create chunk index
	ci := a.makeChunkIndex()

	delete((*(*ChunkReference)[ci].SpriteData), *key)

	SetChanged(true, ci)
}

func (a *ActorRenderer) UpdateVisibility(visible bool) {
	a.Visible = visible
	a.RemoveSprite(nil)
	a.AddSprite(nil)
}

func (a *ActorRenderer) UpdatePos(x, y, z int) {
	oldKey := a.makeKey()
	a.X = x
	a.Y = y
	a.Z = z
	a.RemoveSprite(oldKey)
	a.AddSprite(a.makeKey())
}
