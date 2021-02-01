package main

import(
	"fmt"
	"./Rendering"
	_"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

)

func handleInput(win *pixelgl.Window) {
	//Handle Input
	//Mouse
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		render.CameraStartMove(win.MousePosition())
	}
	if win.Pressed(pixelgl.MouseButtonLeft) {
		render.CameraContinueMove(win.MousePosition())
	}

	//Set Camera Zoom
	render.CameraZoom(win.MouseScroll().Y)

	if win.JustPressed(pixelgl.KeyP) {
		DBBool=!DBBool
		fmt.Println("change",DBBool)
	}

	//Keys
	//Moving Height cutoff
	if win.JustPressed(pixelgl.KeyEqual) {
		render.IncHeightCutoff(0)
	}
	if win.JustPressed(pixelgl.KeyMinus) {
		render.DecHeightCutoff(WorldDepth)
	}

	//Quitting the game
	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

}
