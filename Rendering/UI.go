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

var fpsGraphIndex = 0
var fpsGraph [64]float64

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

	if DBUIEnabled {
		//Reset to screen coordinates
		win.SetMatrix(pixel.IM)

		//Show FPS
		DrawBackingRect(win)

		//Draw the fps graph
		drawGraph(win, fpsGraph, pixel.R(20, 100, 220, 300))

	}

	//Setup for things used later
	Bounds := win.Bounds()
	basicTxt := text.New(pixel.V(10, Bounds.Max.Y-float64(13*len(DBStrings))-20.0), basicAtlas)
	basicTxt.Color = color.RGBA{255, 255, 255, 255}

	//Draw the text
	if DBUIEnabled {

		if len(DBStrings) < 30 {
			for _, s := range DBStrings {
				fmt.Fprint(basicTxt, s)
			}
		} else {
			fmt.Fprintf(basicTxt, "Too much data would be printed\n")
		}
		DBStrings = nil //[]string{}

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
	//fmt.Println(renderFps)

	fpsGraph[fpsGraphIndex] = fps
	fpsGraphIndex++
	fpsGraphIndex = fpsGraphIndex % 64

}

func drawGraph(win *pixelgl.Window, l [64]float64, r pixel.Rect) {
	imd := imdraw.New(nil)
	imd.Color = color.RGBA{255, 255, 255, 255}

	//imd.Color = colornames.Blueviolet
	imd.EndShape = imdraw.RoundEndShape

	w := (r.Max.X - r.Min.X) / 64
	h := (r.Max.Y - r.Min.Y)

	max := 0.0
	min := 9999999.9

	i := 0
	for i < fpsGraphIndex {
		v := fpsGraph[i]
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}

		imd.Push(pixel.V(float64(i)*w+r.Min.X, (v/300.0)*h+r.Min.Y))
		i++
	}
	for i < 64 {
		v := fpsGraph[i]

		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		imd.Push(pixel.V(float64(i)*w+r.Min.X, (v/300.0)*h+r.Min.Y))
		i++
	}

	imd.Line(1)

	imd.Push(r.Min, r.Max)
	imd.Rectangle(1)

	SendString(fmt.Sprintf("Graph min,max: %.4f %.4f\n", min, max))
	SendString(fmt.Sprintf("Graph min,max2: %.4f %.4f\n", 0.0, 300.0))

	imd.Draw(win)
}
