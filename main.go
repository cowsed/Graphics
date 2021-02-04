package main

import (
	"./People"
	"./Rendering"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/pkg/profile"

	"time"
)

const (
	//WorldWidth is the size of a chunk(x)
	WorldWidth int = 16
	//WorldHeight is the height of a chunk (y)
	WorldHeight int = 16
	//WorldDepth is the depth of a chunk (z)
	WorldDepth int = 160
)

//Globals

//DoVSync controls whether or not to vsync
var DoVSync bool = true

//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap []render.Chunk

//Debug things

//DBBool is a boolean controlled by keys to test random features
var DBBool bool = false

//DBBoolLast is the last state of DBBool to get the rising edge of a change
var DBBoolLast bool = true

func run() {
	//Open the window
	cfg := pixelgl.WindowConfig{
		Title:     "Visualization",
		Bounds:    pixel.R(0, 0, 1000, 600),
		VSync:     DoVSync,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Main Loop
	for !win.Closed() {
		last := time.Now() //For FPS Calculations

		//Testing
		render.SendString(fmt.Sprintf("Debug Boolean: %t\n", DBBool))
		if DBBool {
			render.SetAllChanged(true)
		}
		//Input Handling
		inputLast := time.Now()
		handleInput(win)
		inputDt := time.Since(inputLast).Seconds()
		render.SendString(fmt.Sprintf("Input Time(ms): %f\n", 1000*inputDt/60.0))

		/*
			        if DBBoolLast!=DBBool{
			        	DoVSync=DBBool
						win.SetVSync(DoVSync)
				    }
					DBBoolLast=DBBool
		*/

		//personRenderer.RemoveSprite()

		render.SendString(fmt.Sprintf("Vsync: %t\n", DoVSync))
		//Render the world
		drawStart := time.Now()
		render.Render(win, WorldWidth, WorldHeight, WorldDepth)
		//Timing things
		drawDt := time.Since(drawStart).Seconds()
		render.SendString(fmt.Sprintf("Full Self Render Time(ms): %.2f\n", 1000*drawDt))

		upTime := time.Now()
		//Update Window
		win.Update()
		upTimeEnd := time.Since(upTime).Seconds()
		render.SendString(fmt.Sprintf("Window Update Time(ms): %.2f\n", upTimeEnd*1000))

		//Fps calculations
		dt := time.Since(last).Seconds()
		render.SendString(fmt.Sprintf("All Time(ms): %.2f\n", dt*1000))
		//last = time.Now()
		render.SendFPS(1 / dt)

	}
}

func main() {
	//defer profile.Start().Stop()

	//Load Renderer
	render.InitRender()

	chunkData := GenMap2(WorldWidth, WorldHeight, WorldDepth)
	//Make the World. RN a bit hacky (very hacky)
	for i := 0; i < 25; i++ {
		//fmt.Println("Making chunks")
		SpriteData := make(map[[3]int]*render.ActorRenderer)
		chunk := render.Chunk{MaxHeight: WorldDepth, WorldData: &chunkData, SpriteData: &SpriteData, W: WorldWidth, H: WorldHeight, D: WorldDepth}
		chunk.CalculateMax()

		WorldMap = append(WorldMap, chunk)
	}

	//RN this is off because really the renderer should have the chunks but since rn the chuinkmap and the world map are the smae this is what were stuck with
	render.ChunkReference = &WorldMap
	render.SetAllChanged(true)

	//Add test sprite to test sprite rendering
	personRenderer := &render.ActorRenderer{Sheet: nil, FrameIndex: 120, ChunkX: 0, ChunkY: 0}
	person := people.Person{Name: "Timothy", X: 0, Y: 1, Z: 12, Renderer: personRenderer}
	person.UpdateRenderAll(true)
	personRenderer.AddSprite(nil)
	fmt.Println(person)

	personRenderer2 := &render.ActorRenderer{Sheet: nil, FrameIndex: 87, ChunkX: 0, ChunkY: 1}
	person2 := people.Person{Name: "Timothy2", X: 0, Y: 0, Z: 12, Renderer: personRenderer2}
	person2.UpdateRenderAll(true)
	personRenderer2.AddSprite(nil) //Adding sprite and making it calculate its position
	fmt.Println(person)

	pixelgl.Run(run)
}
