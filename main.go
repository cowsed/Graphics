package main

import (
	"./Rendering"
	"./People"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"



	"time"
)
const(
    //WorldWidth is the size of the world (x)
    WorldWidth int =16
     //WorldHeight is the height of generated world (y)
    WorldHeight int =16
    //WorldDepth is the depth of generated world z
    WorldDepth int = 16
)
//Globals


//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap [][][]int

//Debug things
var DBBool bool = false
var DBBoolLast bool = false



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
	fmt.Println(person)

	last := time.Now() //For FPS Calculations
	//Main Loop
	for !win.Closed() {

		//Fps calculations
		dt := time.Since(last).Seconds()
		last = time.Now()
		render.SendFPS(1/dt)


		//Input Handling
		handleInput(win)

        if DBBoolLast!=DBBool{
			if DBBool{
	        	person.X=0
	        } else {person.X=1}
	        person.UpdateRenderAll(true)
	    }
		DBBoolLast=DBBool
		//personRenderer.RemoveSprite()

		//Render the world
		render.Render(win, &WorldMap)

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

