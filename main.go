package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"time"
	"./Rendering"
	
)

//Globalness
var VoidColor color.RGBA = colornames.Skyblue
var WorldMap [][][]int

//Camera Globals
var (
	camPos       = pixel.V(-500,-300)
	oldCamPos    = pixel.ZV
	
	heightCutoff = 0
	mouseStart   = pixel.ZV
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2
)

var mapUpdated = true

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Visualization",
		Bounds: pixel.R(0, 0, 1000, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()

	//Main Loop
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(0, 0), basicAtlas)

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
		if win.JustPressed(pixelgl.KeyEqual) {
			heightCutoff--;
			mapUpdated=true
		}		
		if win.JustPressed(pixelgl.KeyMinus) {
			heightCutoff++;
			mapUpdated=true
		}


		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(pixel.ZV.Sub(camPos))
		oppCam:=pixel.IM.Moved(camPos).Scaled(camPos, 1/camZoom)
		win.SetMatrix(cam)

		win.Clear(VoidColor)

		render.DrawMap(win, WorldMap, mapUpdated, heightCutoff, basicTxt)
		mapUpdated=false
			
		//render.DrawLines(len(WorldMap[0][0]),len(WorldMap[0]),4,win)
		render.DrawSprites(win)

		baseMx:=oppCam//pixel.IM.Moved(win.Bounds().Center().Add(camPos))//Scaled(camPos,1/camZoom).Moved(pixel.ZV)
		fmt.Fprintf(basicTxt,"FPS: %d", int(1/dt))

		basicTxt.Draw(win,  baseMx.Scaled(camPos,2))//.Moved(pixel.V(-400,-300)))

		//Update Window
		win.Update()
	}
}

func main() {
	//Load Background Images
	render.BackgroundInit()
	render.SpritesInit()
	WorldMap = GenMap2(64,64,6)
	//Begin
	pixelgl.Run(run)
}
