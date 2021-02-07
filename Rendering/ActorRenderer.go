package render

import (
	"github.com/faiface/pixel"
)

//This pattern was stolen from a book
type ActorRenderer struct {
	Sheet          *pixel.Picture
	FrameIndex     int
	ChunkX, ChunkY int
	X, Y, Z        int
	Visible        bool
	
	Sprite *pixel.Sprite //This is the actual part that gets referenced. Other parts are for holding excess for animations and such
	
}


func (a *ActorRenderer) Init(){
	a.Sheet=&spriteSheet
	(*a).Sprite=pixel.NewSprite(*a.Sheet, sheetFrames[a.FrameIndex])

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


//Used For animations
//Sets the sprite renderer to use the sprite specified in the actor renderer - called when drawing
func (a *ActorRenderer) UpdateSprite(){
	
	a.Sprite.Set(*a.Sheet, sheetFrames[a.FrameIndex])
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
