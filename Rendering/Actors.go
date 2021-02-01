package render

import (
	_ "fmt"
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
	
}

func (a ActorRenderer) RenderSprite(x,y,z int) {
	//Add the current sprite to the pool to be rendered
	key := string(int(x)) + "" + string(int(y)) + "," + string(int(z))
	SpritesToDraw[key] = &a
	SetChanged(true)
}
