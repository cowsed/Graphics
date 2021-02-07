package render

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

var basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

var fpsSum = 0.0
var fpsCount = 0
var renderFps float64
var maxFps = 0.0
var minFps = 40000.0
var DBStrings []string

var DBUIEnabled = true

func RenderUI(win *pixelgl.Window) {

	//Shows the DB calc of game space
	SendString(CalculateGamePosition(win, win.MousePosition()).String() + "\n")

	//Calculate the game positions 
	gpStart := time.Now()
	gp := CalculateGamePosition(win, win.MousePosition())
	SendString(fmt.Sprintf("CalcMousePosition(ms): %f\n", time.Since(gpStart).Seconds()*1000))

	isox, isoy := isoToWorldCoords(gp)
	SendString(fmt.Sprintf("In World (Base): (%d,%d)\n", isox, isoy))

	//Setup for things used later
	Bounds := win.Bounds()
	basicTxt := text.New(pixel.V(10, Bounds.Max.Y-float64(13*len(DBStrings))-20.0), basicAtlas)
	//Reset to screen coordinates
	win.SetMatrix(pixel.IM)
	

	//Draw the text
	if DBUIEnabled {
		//Show FPS
		DrawBackingRect(win)

		if len(DBStrings) < 30 {
			for _, s := range DBStrings {
				fmt.Fprint(basicTxt, s)
			}
		} else {
			fmt.Fprintf(basicTxt, "Too much data would be printed\n")
		}
		DBStrings = nil//[]string{}

		basicTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.)) // baseMx.Scaled(camPos, 2))
	}
}

//Draw a slightly transparent rectangke for the back if the Debug UI
func DrawBackingRect(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	imd.Color = color.RGBA{0, 0, 0, 150}
	imd.Push(pixel.V(0, 0), pixel.V(300, 700))
	imd.Rectangle(0)
	imd.Draw(win)
}

//Toggle the Debug UI
func ToggleUI() {
	DBUIEnabled = !DBUIEnabled
}


//Send an arbitrary string to the Debug UI
func SendString(s string) {
	//Send an arbitrary string to the UI to write
	DBStrings = append(DBStrings, s)
}

//Send and calculate fps based stats
func SendFPS(fps float64) { //Maybe turn this into a send data function later to send all db data
	renderFps = fps
	minFps = math.Min(renderFps, minFps)
	maxFps = math.Max(renderFps, maxFps)
	fpsSum += fps
	fpsCount++
	SendString(fmt.Sprintf("FPS: %d\n", int(renderFps)))
	SendString(fmt.Sprintf("Avg FPS: %d\n", int(fpsSum/float64(fpsCount))))
	SendString(fmt.Sprintf("Min FPS: %d\n", int(minFps)))
	SendString(fmt.Sprintf("Max FPS: %d\n", int(maxFps)))
}
