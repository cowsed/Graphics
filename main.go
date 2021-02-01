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
const(
    //WorldWidth is the size of the world (x)
    WorldWidth int =16
     //WorldHeight is the height of generated world (y)
    WorldHeight int =16
    //WorldDepth is the depth of generated world z
    WorldDepth int = 20
)
//Globals

//VoidColor defines the color of the background
var VoidColor color.RGBA = colornames.Skyblue
//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap [][][]int

//Debug things
var DBBool bool = false
var DBBoolLast bool = false
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
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Add test sprite to test sprite rendering
    personRenderer:=&render.ActorRenderer{Sheet: nil, FrameIndex: 120}
	person := people.Person{Name: "Timothy",X: 0,Y: 1,Z: 12,Renderer: personRenderer}
    person.UpdateRenderAll(true)
	personRenderer.AddSprite("")


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

        if DBBoolLast!=DBBool{
			if DBBool{
	        	person.X=0
	        } else {person.X=1}
	        person.UpdateRenderAll(DBBool)
	    }
		DBBoolLast=DBBool
		//personRenderer.RemoveSprite()


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
		fmt.Fprintf(basicTxt, "FPS: %d\n", int(1/dt))
		fmt.Fprintf(basicTxt, "Height Cutoff: %d", WorldDepth-heightCutoff)
		basicTxt.Draw(win, baseMx.Scaled(camPos, 2))

		//Update Window
		win.Update()
	}
}

func main() {

	//Load Renderer
	render.InitRender()

	WorldMap = GenMap2(WorldWidth, WorldHeight, WorldDepth)
	//Begin
	pixelgl.Run(run)
}

