package main

import (
	"encoding/json"
	"fmt"
	"./People"
	"./Rendering"
)

type Location struct {
	X int //json:posx
	Y int //json:posy

	W int //json:width
	H int //json:height
	D int //json:depth

	Actors []*people.Actor //json:actors

	Props []string //json:props

	Environment *[][][]int //json:environment
}


//Add an actor to the location - also update it in the renderer - should probably be a different function
func (l *Location) AddActor(a *people.Actor){
	a.Renderer.ChunkX=l.X
	a.Renderer.ChunkY=l.Y

	l.Actors=append(l.Actors, a)
	a.Renderer.Init()
	a.UpdateRenderAll(true)
}


//Create the render version
func (l *Location) MakeRenderChunk() render.Chunk {

	//Create
	chunk := render.Chunk{MaxHeight: -1, WorldData: l.Environment, W: l.W, H: l.H, D: l.D}

	//Add all the actors
	for _, actor := range l.Actors {

		actor.Renderer.Init()
	}

	//Set Attributes
	chunk.CalculateMax()
	chunk.SetChanged(true)

	return chunk

}

//Jsoning
func (l *Location) Marshal() string {

	result, _ := json.Marshal(l)
	fmt.Println(string(result))
	
	return string(result)
}
