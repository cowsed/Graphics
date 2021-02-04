package render

import (
	_"fmt"
)

//A type to hold the chunk data in one place
type Chunk struct {
	changed    bool
	MaxHeight  int
	WorldData  *[][][]int                 //Holds the tile materials
	SpriteData *map[[3]int]*ActorRenderer //Holds the character sprites to be drawn
	W, H, D    int
}

//Consider adding a way that only runs through all 3Ds and saves the tiles that are actually drawn and their position as well as if they were visible or not
//Then use the dirty pattern to mark it to recalculate if a tile or a sprite changes

//unfortunaetly i dont think this will work any better than just checking changed before rendering when sprites are moving
//sike no i think it would
//As long as you assume sprites dont cover the entire block (and even then it doesnt have  z aissue itll just be rendering more than necessary)
//You can run through the list of the tiles that are to be drawn and then keep a list of sprites and interject them when necessary
//THe interjecting part may be a bit difficult because youll have to hold the keys but that shouldnt be tooo? hard
//Only need to say its dirty if a tile has changed,
//if the height cutoff has changed one can think you can just cutoff the top ones but that would mess up the hidden style drawing

//Calculators
//Calculate the maximum non air block
func (c *Chunk) CalculateMax() {
	//Run down from top and stop once a non air block is hit
	MaxHeight := -1

	for z := c.D - 1; z >= 0; z-- {
		for y := 0; y < c.H; y++ {
			for x := 0; x < c.W; x++ {
				//If theres a high tile
				if (*c.WorldData)[z][y][x] != 0 {
					MaxHeight = z
					break
				}
				//If theres a high sprite
				if _, ok := (*c.SpriteData)[[3]int{x, y, z}]; ok {
					MaxHeight = z
					break

				}
			}
			if MaxHeight != -1 {
				break
			}
		}
		if MaxHeight != -1 {
			break
		}
	}

	(*c).MaxHeight = MaxHeight+1
}

//Setters
//Set if it has changed
func (c *Chunk) SetChanged(changed bool) {
	c.changed = changed
}

//Getters
func (c *Chunk) GetChanged() bool {
	return c.changed
}
