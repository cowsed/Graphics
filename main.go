package main

import (
	"fmt"
	"time"

	"./People"
	"./Rendering"
	"./Materials"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"

)

//Debug things

//DOProfile is a debug variable saying whether or not to make a profile when it runs
const DOProfile = false

//DBBool is a boolean controlled by keys to test random features
var DBBool bool = false

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
var DoVSync bool = false

//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap []render.Chunk


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

		//Debug Logging of VSync
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
	//If Specified log to a profile
	if DOProfile{
		defer profile.Start(profile.ProfilePath(".")).Stop()
	}
	
	
	
	
	//Make the World. RN a bit hacky (very hacky)
	for i := 0; i < render.NumChunks; i++ {
		//fmt.Println("Making chunks")

		chunky := 5 - i/5
		chunkx := i % 5
		chunkData := GenMap3(WorldWidth, WorldHeight, WorldDepth, chunkx, chunky)

		chunk := render.Chunk{MaxHeight: WorldDepth, WorldData: &chunkData, W: WorldWidth, H: WorldHeight, D: WorldDepth}
		chunk.CalculateMax()

		WorldMap = append(WorldMap, chunk)
	}

	//RN this is off because really the renderer should have the chunks but since rn the chuinkmap and the world map are the smae this is what were stuck with
	render.ChunkReference = &WorldMap
	render.SetAllChanged(true)

	//Load Renderer
	render.InitRender()

	//Add test sprite to test sprite rendering
	//Initializing things like this is rather wasteful as it creates and recalculates many things many times
	personRenderer := &render.ActorRenderer{Sheet: nil, FrameIndex: materials.ROCK_WALL_V_1-1, ChunkX: 0, ChunkY: 0}
	person := people.Person{Name: "Timothy", X: 0, Y: 0, Z: 10, Renderer: personRenderer}
	personRenderer.Init()
	
	fmt.Println(person)

	person.UpdateRenderAll(true)
	
	fmt.Println(person)

	personRenderer2 := &render.ActorRenderer{Sheet: nil, FrameIndex: materials.ROCK_WALL_V_1-1, ChunkX: 0, ChunkY:0}
	person2 := people.Person{Name: "Timothy2", X: 0, Y: 1, Z: 10, Renderer: personRenderer2}
	personRenderer2.Init()

	fmt.Println(person2)

	person2.UpdateRenderAll(true)
	fmt.Println(person2)

	pixelgl.Run(run)
}
