package render

import (
	"fmt"
	//"math/rand"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_"github.com/faiface/pixel/text"
	_"image"
	_ "image/png"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"


)

var actorSheet pixel.Picture
var actorFrames []pixel.Rect
var person Sprite


func SpritesInit(){
	actorSheet,actorFrames=loadSheet("Assets/sprites2.png",32,64)
	fmt.Println("Loaded sprites")
	person=Sprite{&actorSheet,&actorFrames, 12, 0,0,5}
}

//Oh boy a flywieght pattern
type Sprite struct{
	sheet *pixel.Picture
	frames *[]pixel.Rect
	index int
	//Array Space
	x,y,z int
}

func DrawSprites(win *pixelgl.Window ){//, sprites []Sprite
	//NEED TO DO DEPTH TESTING STUFF

	sprites:=[]Sprite{person}
	for _,s := range sprites{
		//place=toIsoCoords(s.x,s.y,s.z)
		mx:=pixel.IM.Moved(toIsoCoords(s.x,s.y,s.z)).Moved(pixel.V(-4,10.0))
		sprite := pixel.NewSprite(*s.sheet, (*s.frames)[s.index])
		sprite.Draw(win, mx)
	}

}

func DrawLines(w,h,d int, win *pixelgl.Window){
	lineWidth:=2.0
	imd := imdraw.New(nil)


		for x:=0; x<w; x++{
			imd.Color = colornames.Black
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(  toIsoCoords(x-1,0,d), toIsoCoords(x-1,h,d)   )
			imd.EndShape = imdraw.SharpEndShape
			imd.Line(lineWidth)
		}
		
		for y:=0; y<h; y++{
			imd.Color = colornames.Black
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(  toIsoCoords(-1,y,d), toIsoCoords(w-1,y,d)   )
			imd.EndShape = imdraw.SharpEndShape
			imd.Line(lineWidth)		
		}
	imd.Draw(win)

}


func toIsoCoords(x,y,z int) pixel.Vec {
	return pixel.V(
		float64(32.0*(x))+float64(32*y),
		float64(16.0*(y-x)+z*32))
}
