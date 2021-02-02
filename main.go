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
    WorldWidth int =32
     //WorldHeight is the height of generated world (y)
    WorldHeight int =32
    //WorldDepth is the depth of generated world z
    WorldDepth int = 16
)
//Globals


//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap []*[][][]int

//Debug things

//DBBool is a boolean controlled by keys to test random features
var DBBool bool = false
//DBBoolLast is the last state of DBBool to get the rising edge of a change
var DBBoolLast bool = true



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
    personRenderer:=&render.ActorRenderer{Sheet: nil, FrameIndex: 120, ChunkX: 0, ChunkY: 0, }
	person := people.Person{Name: "Timothy",X: 0,Y: 1,Z: 12,Renderer: personRenderer}
    person.UpdateRenderAll(true)
	personRenderer.AddSprite("")
	fmt.Println(person)

    personRenderer2:=&render.ActorRenderer{Sheet: nil, FrameIndex: 87, ChunkX: 0, ChunkY: 1, }
	person2 := people.Person{Name: "Timothy2",X: 0,Y: 1,Z: 12,Renderer: personRenderer2}
    person2.UpdateRenderAll(true)
	personRenderer2.AddSprite("")
	fmt.Println(person)



	//Main Loop
	for !win.Closed() {
	last := time.Now() //For FPS Calculations


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
		render.Render(win, WorldMap,WorldWidth,WorldHeight)

		//Update Window
		win.Update()
	
		//Fps calculations
		dt := time.Since(last).Seconds()
		last = time.Now()
		render.SendFPS(1/dt)

	}
}

func main() {

	//Load Renderer
	render.InitRender()

	//This is a really hacky way to create a 3x3 grid cuz the generator isnt there yet
	chunk:= GenMap2(WorldWidth, WorldHeight, WorldDepth)
	WorldMap=[]*[][][]int{&chunk,&chunk,&chunk,&chunk,&chunk,&chunk,&chunk,&chunk,&chunk}
	//Begin
	pixelgl.Run(run)
}

