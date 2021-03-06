package main

import (
	"fmt"
	"time"

	"./Config"
	"./Materials"
	"./Rendering"
	"./World"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//Further Profiling
import _ "net/http/pprof"
import "net/http"

//Debug things

//DOProfile is a debug variable saying whether or not to make a profile when it runs
const DOProfile = false

//DBBool is a boolean controlled by keys to test random features
var DBBool bool = false

//Globals

//DoVSync controls whether or not to vsync
var DoVSync bool = false

//WorldMap holds the world grid that holds the world (may soon be changed to a more memory friendly version
var WorldMap []Location

var RenderChunkMap []render.Chunk

//Testing
//Peoples
var person people.Actor

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
			person.Update()
			//render.SetAllChanged(true)
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
		render.Render(win, config.ChunkWidth, config.ChunkHeight, config.ChunkDepth)

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
	if DOProfile {
		//defer profile.Start(profile.ProfilePath(".")).Stop()

		go func() {
			fmt.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	materials.LoadSprites("Materials/tiles.csv", "Assets/Custom3.png")

	//Make the World. RN a bit hacky (very hacky)
	for i := 0; i < render.NumChunks; i++ {
		//fmt.Println("Making chunks")

		chunky := i / 5
		chunkx := i % 5
		chunkData := GenMap3(config.ChunkWidth, config.ChunkHeight, config.ChunkDepth, chunkx, chunky)

		chunk := Location{X: chunkx, Y: chunky, W: config.ChunkWidth, H: config.ChunkHeight, D: config.ChunkDepth, Actors: []*people.Actor{}, Props: []string{}, Environment: &chunkData}

		//Testing marshalling
		s := chunk.Marshal()

		newChunk := Location{}
		newChunk.Unmarshal([]byte(s))

		WorldMap = append(WorldMap, newChunk)
		RenderChunkMap = append(RenderChunkMap, chunk.MakeRenderChunk())
	}

	//RN this is off because really the renderer should have the chunks but since rn the chuinkmap and the world map are the smae this is what were stuck with
	render.ChunkReference = &RenderChunkMap

	//Load Renderer
	render.InitRender()

	//Add test sprite to test sprite rendering
	//Initializing things like this is rather wasteful as it creates and recalculates many things many times
	personRenderer := &render.ActorRenderer{Sheet: nil, FrameIndex: materials.Sprites["ROCK_WALL_V_1"]}
	person = people.Actor{Name: "Timothy", X: 0, Y: 0, Z: 10, Renderer: personRenderer}
	WorldMap[0].AddActor(&person)

	fmt.Println(person)

	personRenderer2 := &render.ActorRenderer{Sheet: nil, FrameIndex: materials.Sprites["GRASS_LESS_2"]}
	person2 := people.Actor{Name: "Timothy2", X: 0, Y: 1, Z: 10, Renderer: personRenderer2}

	//person2.AddToChunk(&person2)
	WorldMap[0].AddActor(&person2)

	fmt.Println(person2)

	pixelgl.Run(run)
}
