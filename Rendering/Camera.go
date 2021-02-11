package render

import (
	"github.com/faiface/pixel"
	"math"
)

//Camera Globals
var (
	camPos       = pixel.V(0, 0)
	oldCamPos    = pixel.ZV
	mouseStart   = pixel.ZV
	camZoom      = 1.0
	camZoomSpeed = 1.2

	heightCutoff = 1
)

//Set the camera zoom
func CameraZoom(value float64) {
	camZoom *= math.Pow(camZoomSpeed, value)
}

//Begin a camera move
func CameraStartMove(mousePos pixel.Vec) {
	oldCamPos = camPos
	mouseStart = mousePos
}

//Continue the camera move started by CameraStartMove
func CameraContinueMove(mousePos pixel.Vec) {
	camPos = oldCamPos.Add(mouseStart.Sub(mousePos).Scaled(1.0 / camZoom))
}

//Increment the Height Cutoff
func IncHeightCutoff(min int) {
	if heightCutoff > min {
		heightCutoff--
		SetAllDirty(true)
		SetAllChanged(true)

	}
}

//Decrement the Height Cutoff
func DecHeightCutoff(max int) {
	if heightCutoff < max-1 {
		heightCutoff++
		SetAllDirty(true)
		SetAllChanged(true)

	}
}
