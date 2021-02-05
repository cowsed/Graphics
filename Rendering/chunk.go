package render

import (
	_"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

type Tile struct {
	X, Y, Z   int
	TileIndex int
	//If the tile should be drawn with the "inside" overlay
	Hidden bool
}

//A type to hold the chunk data in one place
type Chunk struct {
	//If the chunk has changed (sprites moved)
	spriteChanged bool

	//If the tiles have changed
	Dirty bool

	//Maximum height to render to
	MaxHeight int

	//Holds the world in the form of tiles
	WorldData *[][][]int //Holds the tile materials

	//Holds the Tiles that have been computed to be visible
	VisibleData []Tile

	//Holds the character sprites to be drawn
	SpriteData *map[[3]int]*ActorRenderer

	//The batch of sprites to be drawn when called upon
	Batch *pixel.Batch

	//Width Height and Depth of the chunk
	W, H, D int
}

//Does all the first time set up stuff necessary to render
func (c *Chunk) Init() {
	//Create the batch
	c.Batch = pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)
	c.Dirty = true
}

//Update the tile map
func (c *Chunk) UpdateTiles(ChunkIndex, cutoffHeight int) {
	if c.Dirty {
		c.FindVisible(ChunkIndex, cutoffHeight)
	}
	//Set it to dirty so it fixes itself on startup
	c.Dirty = false
}

//Render the tile map to the chunk batch
func (c *Chunk) Render(win *pixelgl.Window, chunkX, chunkY int) {
	//Send the relative chunkX and Y coordinates bc they are a state completely made up by the renderer

	//fmt.Println(c.VisibleData)
	//Assemble the batch
	c.RenderToBatch(chunkX, chunkY)
	//Draw to the window
	//c.Batch.Draw(win)
}

//Render the world to the batch
func (c *Chunk) RenderToBatch(chunkX, chunkY int) {

	if c.spriteChanged {

		//Reset the batch
		c.Batch.Clear()
		for _, tile := range c.VisibleData {
			//TODO: At some point in here drag in the Sprites after all thats what changed is all about
			if tile.Hidden { //Set the 'inside' mask
				c.Batch.SetColorMask(color.RGBA{60, 60, 60, 255})
			}
			//sprite sheet need not be passed because its just a global thing
			mx:=pixel.IM.Moved(worldToIsoCoords(tile.X+(chunkX*c.W), tile.Y+(chunkY*c.H), tile.Z))
			sprite := pixel.NewSprite(spriteSheet, sheetFrames[tile.TileIndex-1])
			sprite.Draw(c.Batch, mx)
		
			if tile.Hidden { //Reset the 'inside' mask
				c.Batch.SetColorMask(color.RGBA{255, 255, 255, 255})
			}
		}
		//Reset changed
		c.spriteChanged = false
	} else { //If no changes were made theres no need to redraw the batches so just leave it be
		//return c.Batch
	}

}

//Calculators

//Calculate the pieces that are visible
//Cutoff Height is the toal depth - height cutoff

func (c *Chunk) FindVisible(ChunkIndex, cutoffHeight int) {
	//Reset visibleData slice (i think this is what caused the crash earlier
	//TODO maybe reuse space to make the garbage collector happy as i think the gc frames are occuring causing stuttering
	c.VisibleData = nil
	//Capacity here is totally made up tho hopefully it helps with cacheing and locality
	//TODO make a simple function that would calculate the surface area of a box the size of the chunk or something and use that
	//fmt.Println("Updating Visible Data")
	c.VisibleData = make([]Tile, 0, 400)
	//Run through all of 3d space and include
	for z := 0; z < intMin(c.MaxHeight, c.D-cutoffHeight); z++ {
		for y := c.H - 1; y >= 0; y-- { //Go backwards to go over each other
			for x := 0; x < c.W; x++ {
				tileIndex := (*c.WorldData)[z][y][x]  //Minus 1 because 0 here is air not the 0 sprite
				if tileIndex != 0 {                     //again, now -1 is air
					if visible, inside := CheckVisibility(x, y, z, c.W, c.H, c.D, cutoffHeight, ChunkIndex); visible {
						//If a non-air tile is visible add it to the visible collection

						c.VisibleData = append(c.VisibleData, Tile{X: x, Y: y, Z: z, TileIndex: tileIndex, Hidden: inside})
					}

				}
			}
		}
	}
	//Just in case the caller didnt do this
	//Also needs to redraw everything
	c.spriteChanged = true

}

/*
//Consider adding a way that only runs through all 3Ds and saves the tiles that are actually drawn and their position as well as if they were visible or not
//Then use the dirty pattern to mark it to recalculate if a tile or a sprite changes

unfortunaetly i dont think this will work any better than just checking changed before rendering when sprites are moving
sike no i think it would
As long as you assume sprites dont cover the entire block (and even then it doesnt have  z aissue itll just be rendering more than necessary)
You can run through the list of the tiles that are to be drawn and then keep a list of sprites and interject them when necessary
THe interjecting part may be a bit difficult because youll have to hold the keys but that shouldnt be tooo? hard
Only need to say its dirty if a tile has changed,
if the height cutoff has changed one can think you can just cutoff the top ones but that would mess up the hidden style drawing
*/
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
	(*c).MaxHeight = MaxHeight + 1
}

//Setters

//Set if it has changed
func (c *Chunk) SetChanged(changed bool) {
	c.spriteChanged = changed
}
func (c *Chunk) SetDirty(dirty bool) {
	c.Dirty = dirty
}

//Getters
func (c *Chunk) GetChanged() bool {
	return c.spriteChanged
}
func (c *Chunk) GetDirty() bool {
	return c.Dirty
}
