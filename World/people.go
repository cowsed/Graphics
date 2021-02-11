package people

import (
	"../Materials"
	"../Rendering"
	
)

var AnimationTick = 0

type Actor struct {
	Name     string
	X, Y, Z  int
	Renderer *render.ActorRenderer
}

func (a *Actor) AddToChunk(){

}


func (a *Actor) Update() {
	//Switch sprites for now
	
	AnimationTick++
	
	
	if AnimationTick%12 == 0 {
		if a.Renderer.FrameIndex == materials.Sprites["STONE_1"] {
			a.Renderer.UpdateSprite(materials.Sprites["STONE_2"])
		} else {
			a.Renderer.UpdateSprite(materials.Sprites["STONE_1"])
		}
	}
}

func (a *Actor) UpdateRenderAll(enabled bool) {

	(*a.Renderer).UpdatePos(a.X, a.Y, a.Z)

	//Something weird is happening where passing enabled doesnt update it but just saying false does
	(*a.Renderer).UpdateVisibility(enabled)
}

//Updates the rendered position to be equal to the persons position
func (a *Actor) UpdateRenderPos() { //, sprite offsett
	a.Renderer.UpdatePos(a.X, a.Y, a.Z)
}

func (a *Actor) UpdateVisibility(visible bool) {
	//fmt.Println("To: ",visible)
	(*a.Renderer).UpdateVisibility(visible)
	//fmt.Println("Outcome: ",(*p.Renderer))

}
