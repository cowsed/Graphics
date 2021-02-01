package render

import(
		"fmt"
		"github.com/faiface/pixel"
		"github.com/faiface/pixel/pixelgl"
		"github.com/faiface/pixel/text"
		"golang.org/x/image/font/basicfont"
	)
var basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

var renderFps float64

func RenderUI(win *pixelgl.Window){
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)

	//Reset to screen coordinates
	win.SetMatrix(pixel.IM)
	//Show FPS
	fmt.Fprintf(basicTxt, "FPS: %d\n", int(renderFps))
	//Show Height Cutoff
	fmt.Fprintf(basicTxt, "Height Cutoff: %d",0)// WorldDepth-heightCutoff)

	//Draw the text at 0,0
	basicTxt.Draw(win,pixel.IM.Scaled(pixel.ZV,2))// baseMx.Scaled(camPos, 2))

}


func SendFPS(fps float64){//Maybe turn this into a send data function later to send all db data
	renderFps=fps
}
