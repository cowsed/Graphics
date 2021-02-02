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
	//Controls whether or not to vsync
var DoVSync bool = true

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
		VSync:  DoVSync,
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
		
		//Testing
		//render.SetAllChanged(true)

		//Input Handling
		inputLast:=time.Now()
		handleInput(win)
		inputDt := time.Since(inputLast).Seconds()
		render.SendString(fmt.Sprintf("Input Time(ms): %f\n",1000*inputDt/60.0))

		/*		
        if DBBoolLast!=DBBool{
        	DoVSync=DBBool
			win.SetVSync(DoVSync)
	    }
		DBBoolLast=DBBool
		*/

		//personRenderer.RemoveSprite()

		render.SendString(fmt.Sprintf("Vsync: %t\n",DoVSync))
		//Render the world
		drawStart:=time.Now()
		render.Render(win, WorldMap,WorldWidth,WorldHeight)
		//Timing things
		drawDt := time.Since(drawStart).Seconds()
		render.SendString(fmt.Sprintf("Render Time(ms): %f\n",1000*drawDt))

		upTime:=time.Now() 
		//Update Window
		win.Update()
		upTimeEnd:=time.Since(upTime).Seconds()
		render.SendString(fmt.Sprintf("Update Time(ms): %f\n", upTimeEnd*1000))
	
		//Fps calculations
		dt := time.Since(last).Seconds()
		render.SendString(fmt.Sprintf("All Time(ms): %f\n", dt*1000))
		//last = time.Now()
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

