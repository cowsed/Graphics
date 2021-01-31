package main

import (
	"./Rendering"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
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
	person := render.ActorRenderer{nil, 12, 7, 5, 4}
	person.RenderSprite()

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

	WorldMap = GenMap2(64, 64, 6)
	//Begin
	pixelgl.Run(run)
}

func handleInput(win *pixelgl.Window) {
	//Handle Input
	//Mouse
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		oldCamPos = camPos
		mouseStart = win.MousePosition()
	}
	if win.Pressed(pixelgl.MouseButtonLeft) {

		camPos = oldCamPos.Add(mouseStart.Sub(win.MousePosition()).Scaled(1.0 / camZoom))
	}

	camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

	//Keys
	//
	if win.JustPressed(pixelgl.KeyEqual) {
		if heightCutoff > 0 {
			heightCutoff--
			render.SetChanged(true)
		}
	}
	if win.JustPressed(pixelgl.KeyMinus) {
		if heightCutoff < len(WorldMap)-1 {
			heightCutoff++
			render.SetChanged(true)
		}
	}

	//Quitting the game
	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

}
