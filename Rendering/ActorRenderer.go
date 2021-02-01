package render

import (
	 _"fmt"
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
	Sheet      *pixel.Picture
	FrameIndex int
	X,Y,Z int
	Visible bool
}
func (a ActorRenderer) makeKey() string {
	return string(int(a.X)) + "" + string(int(a.Y)) + "," + string(int(a.Z))
	
}

//Add Sprite adds the sprite to the pool of sprites to be included in the batch
//Key us an arguement because if the position changes the key changes so updating it doesnt work
func (a ActorRenderer) AddSprite(key string) {
	//Add the current sprite to the pool to be rendered
	//Blank key if calling externally
	if key==""{
		key=a.makeKey()
	}
	SpritesToDraw[key] = &a
	SetChanged(true)
}

//Remove the sprite from the pool
func (a ActorRenderer) RemoveSprite(key string){
	//Blank key if calling externally
	if key==""{
		key=a.makeKey()
	}
	delete(SpritesToDraw,key)
	SetChanged(true)
}

func (a *ActorRenderer) UpdateVisibility( visible bool){
	a.Visible=visible
	a.RemoveSprite("")
	a.AddSprite("")
}

func (a *ActorRenderer) UpdatePos(x,y,z int){
	oldKey:=a.makeKey()
	a.X=x
	a.Y=y
	a.Z=z
	a.RemoveSprite(oldKey)
	a.AddSprite( a.makeKey() )
}
