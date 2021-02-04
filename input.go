package main

import (
	"./Rendering"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/png"
	"os"
	"time"
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
		DBBool = !DBBool
		fmt.Println("change", DBBool)
	}
	//Toggle VSync
	if win.JustPressed(pixelgl.KeyV) {
		DoVSync = !DoVSync
		win.SetVSync(DoVSync)
	}

	//Toggle the visibility of the debug UI
	if win.JustPressed(pixelgl.KeyU) {
		render.ToggleUI()
	}

	//Keys
	//Moving Height cutoff
	if win.JustPressed(pixelgl.KeyEqual) {
		render.IncHeightCutoff(0)
	}
	if win.JustPressed(pixelgl.KeyMinus) {
		render.DecHeightCutoff(WorldDepth)
	}

	if win.JustPressed(pixelgl.KeySlash) {
		fmt.Println("TakingScreenshot")
		fmt.Println("taking screenshot...")

		f, err := os.Create(fmt.Sprint("Screenshots/Screenshot-", time.Now().Format("1-2-3:4:5")))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		img := pixel.PictureDataFromPicture(win)
		png.Encode(f, img.Image())

		fmt.Println("done")
	}

	//Quitting the game
	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

}
