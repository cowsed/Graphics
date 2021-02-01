package main

import(
	"fmt"
	"./Rendering"
	_"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

func handleInput(win *pixelgl.Window) {
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


	if win.JustPressed(pixelgl.KeyP) {
		DBBool=!DBBool
		fmt.Println("change",DBBool)
	}

	//Keys
	//Moving Height cutoff
	if win.JustPressed(pixelgl.KeyEqual) {
		if heightCutoff > 0 {
			heightCutoff--
			render.SetChanged(true)
		}
	}
	if win.JustPressed(pixelgl.KeyMinus) {
		if heightCutoff < len(WorldMap)-1 {
			heightCutoff++
			render.SetChanged(true)
		}
	}

	//Quitting the game
	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

}
