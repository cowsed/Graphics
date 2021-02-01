package main

import (
	"./Rendering"
	"./People"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"

	"time"
)


//Globalness
var VoidColor color.RGBA = colornames.Skyblue
var WorldMap [][][]int

//Camera Globals
var (
	camPos    = pixel.V(-500, -300)
	oldCamPos = pixel.ZV

	heightCutoff = 0
	mouseStart   = pixel.ZV
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2
)

var mapUpdated = true

func run() {
	//Open the window
	cfg := pixelgl.WindowConfig{
		Title:  "Visualization",
		Bounds: pixel.R(0, 0, 1000, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Add test sprite to test sprite rendering
	person := people.Person{"Timothy",0,0,1,&render.ActorRenderer{nil, 100}}
	person.Render()

	last := time.Now() //For FPS Calculations
	//Main Loop
	for !win.Closed() {
		//Fps calculations
		dt := time.Since(last).Seconds()
		last = time.Now()

		//Text Setup
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(0, 0), basicAtlas)

		//Input Handling
		handleInput(win)

		//Calculate camera positioning and UI positioning
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(pixel.ZV.Sub(camPos))
		oppCam := pixel.IM.Moved(camPos).Scaled(camPos, 1/camZoom)
		win.SetMatrix(cam)

		//Clear the Window to prepare for drawing
		win.Clear(VoidColor)

		render.RenderAll(win, &WorldMap, heightCutoff)
		//render.DrawLines(len(WorldMap[0][0]),len(WorldMap[0]),4,win)

		mapUpdated = false

		//UI Stuff
		baseMx := oppCam
		fmt.Fprintf(basicTxt, "FPS: %d", int(1/dt))
		basicTxt.Draw(win, baseMx.Scaled(camPos, 2))

		//Update Window
		win.Update()
	}
}

func main() {

	//Load Renderer
	render.InitRender()

	WorldMap = GenMap2(140, 140, 8)
	//Begin
	pixelgl.Run(run)
}

