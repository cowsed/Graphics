package render

import(
		"fmt"
		"math"
		"github.com/faiface/pixel"
		"github.com/faiface/pixel/pixelgl"
		"github.com/faiface/pixel/text"
		"golang.org/x/image/font/basicfont"
	)
var basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)


var fpsSum =0.0
var fpsCount=0
var renderFps float64
var maxFps = 0.0
var minFps = 40000.0
var DBStrings []string

func RenderUI(win *pixelgl.Window){
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)

	//Reset to screen coordinates
	win.SetMatrix(pixel.IM)
	//Show FPS
	fmt.Fprintf(basicTxt, "FPS: %d\n", int(renderFps))
	fmt.Fprintf(basicTxt, "Avg FPS: %d\n", int(fpsSum/float64(fpsCount)))
	fmt.Fprintf(basicTxt, "Min FPS: %d\n", int(minFps))
	fmt.Fprintf(basicTxt, "Max FPS: %d\n", int(maxFps))
	
	if len(DBStrings)<10{
		for _,s := range DBStrings{
			fmt.Fprintf(basicTxt, s)
		}
	} else {
		fmt.Fprintf(basicTxt, "Too much data would be printed")		
	}
	DBStrings=[]string{}
	//Show Height Cutoff
	fmt.Fprintf(basicTxt, "Height Cutoff: %d",16-heightCutoff)

	//Draw the text at 0,0
	basicTxt.Draw(win,pixel.IM.Scaled(pixel.ZV,2))// baseMx.Scaled(camPos, 2))

}

func SendString(s string){
	//Send an arbitrary string to the UI to write
	DBStrings=append(DBStrings,s)
}

func SendFPS(fps float64){//Maybe turn this into a send data function later to send all db data
	renderFps=fps
	minFps=math.Min(renderFps,minFps)
	maxFps=math.Max(renderFps,maxFps)
	fpsSum+=fps
	fpsCount++
}
