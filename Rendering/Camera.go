package render

import(
	"math"
	"github.com/faiface/pixel"
)


//Camera Globals
var (
	camPos    = pixel.V(800, 800)
	oldCamPos = pixel.ZV
	mouseStart   = pixel.ZV
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2
	
	heightCutoff = 0
)

func CameraZoom(value float64){
	camZoom *= math.Pow(camZoomSpeed, value)
}

func CameraStartMove(mousePos pixel.Vec){
	oldCamPos = camPos
	mouseStart = mousePos
}
func CameraContinueMove(mousePos pixel.Vec){
	camPos = oldCamPos.Add(mouseStart.Sub(mousePos).Scaled(1.0 / camZoom))
}

func IncHeightCutoff(min int){
	if heightCutoff > min {
		heightCutoff--
		SetAllChanged(true)
	}
}

func DecHeightCutoff(max int){
	if heightCutoff < max-1 {
		heightCutoff++
		SetAllChanged(true)
	}
}
