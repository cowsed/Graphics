package render

import (
	_"fmt"
	//"math/rand"
	"github.com/faiface/pixel"
	_"github.com/faiface/pixel/imdraw"
	_"github.com/faiface/pixel/pixelgl"
	_ "github.com/faiface/pixel/text"
	_"golang.org/x/image/colornames"
	_ "image"
	_ "image/png"
)

//This pattern was stolen from a book
type ActorRenderer struct {
	Sheet      *pixel.Picture
	FrameIndex int
	X,Y,Z int
}

func (a ActorRenderer) RenderSprite() {
	//Add the current sprite to the pool to be rendered
	key:=string(int(a.X))+""+string(int(a.Y))+","+string(int(a.Z))
	SpritesToDraw[key]=&a
}



