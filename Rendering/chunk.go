package render

import (
	"fmt"
	"image/color"
	"sort"

	"../Config"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Tile struct {
	X, Y, Z   int
	TileIndex int
	Sprite    *pixel.Sprite
	//If the tile should be drawn with the "inside" overlay
	Hidden bool
}

//A type to hold the chunk data in one place
type Chunk struct {
	//If the chunk has changed (sprites moved)
	spriteDirty bool

	//If the tiles have changed
	tileDirty bool

	//Maximum height to render to
	MaxHeight int

	//Holds the world in the form of tiles
	WorldData *[config.ChunkDepth][config.ChunkHeight][config.ChunkWidth]int //Holds the tile materials

	//Holds the Tiles that have been computed to be visible
	VisibleData []Tile

	//Holds the character sprites to be drawn
	SpriteDataOrdered []*ActorRenderer

	//The batch of sprites to be drawn when called upon
	Batch *pixel.Batch

	//Width Height and Depth of the chunk
	W, H, D int
}

//Does all the first time set up stuff necessary to render
func (c *Chunk) Init() {
	//Create the batch
	c.Batch = pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)

	//Tell the system to recalculate itself
	c.spriteDirty = true
	c.tileDirty = true
}

//Update the tile map
func (c *Chunk) UpdateTiles(ChunkIndex, cutoffHeight int) {
	if c.tileDirty {
		c.FindVisible(ChunkIndex, cutoffHeight)

		c.SetDirty(false)
	}

}

//Render the tile map to the chunk batch
func (c *Chunk) Render(win *pixelgl.Window, chunkX, chunkY int) {
	//Send the relative chunkX and Y coordinates bc they are a state completely made up by the renderer

	//Assemble the batch
	c.RenderToBatch(chunkX-1, chunkY)
	//Draw to the window
	c.Batch.Draw(win)
}

//Render the world to the batch
//TODO: Look into the cache locality of this solution (a sprite interface may work better but it also may not)
func (c *Chunk) RenderToBatch(chunkX, chunkY int) {
	if c.spriteDirty { //If no changes were made theres no need to redraw the batches so just leave it be

		//Reset the batch
		c.Batch.Clear()

		//Resort the sprites based on position
		c.sortARs(c.SpriteDataOrdered)

		//Similar to the i=y*w+x type thing just 3d
		spriteIndex := 0 //The index of the ActorRenderers

		for _, tile := range c.VisibleData {

			tileScore := c.calcScore(tile.X, tile.Y, tile.Z)
			//While There are still strites to draw in the layer, draw them

			if len(c.SpriteDataOrdered) > 0 && spriteIndex < len(c.SpriteDataOrdered)-1 { //If there are any sprites and we've not run out of sprites
				ar := c.SpriteDataOrdered[spriteIndex] //Next ActorRenderer

				for c.calcScore(ar.X, ar.Y, ar.Z) < tileScore {
					//Draw the sprite
					mx := pixel.IM.Moved(worldToIsoCoords(ar.X+(chunkX*c.W), ar.Y+(chunkY*c.H), ar.Z))
					(*ar).UpdateSpriteImage() //pixel.NewSprite(spriteSheet, sheetFrames[(*ar).FrameIndex])
					ar.Sprite.Draw(c.Batch, mx)

					if spriteIndex < len(c.SpriteDataOrdered)-1 {
						spriteIndex++
						ar = c.SpriteDataOrdered[spriteIndex] //Next ActorRenderer
					} else {
						break
					}
				}
			}

			if tile.Hidden { //Set the 'inside' mask
				c.Batch.SetColorMask(color.RGBA{60, 60, 60, 255})
			}

			mx := pixel.IM.Moved(worldToIsoCoords(tile.X+(chunkX*c.W), tile.Y+(chunkY*c.H), tile.Z))

			//If mx is within camera view

			tile.Sprite.Draw(c.Batch, mx)

			if tile.Hidden { //Reset the 'inside' mask
				c.Batch.SetColorMask(color.RGBA{255, 255, 255, 255})
			}

		}
		//If finished and theres still some sprites left over, draw them

		//Draw the sprite
		for i := spriteIndex; i < len(c.SpriteDataOrdered); i++ {
			ar := c.SpriteDataOrdered[i]

			mx := pixel.IM.Moved(worldToIsoCoords(ar.X+(chunkX*c.W), ar.Y+(chunkY*c.H), ar.Z))
			(*ar).UpdateSpriteImage() //pixel.NewSprite(spriteSheet, sheetFrames[(*ar).FrameIndex])
			ar.Sprite.Draw(c.Batch, mx)
		}

		//Reset changed
		c.spriteDirty = false
	}
}

//Caluclates the single value for position
func (c *Chunk) calcScore(x, y, z int) int {
	return (z * c.W * c.H) + (((c.W - 1) - y) * c.W) + x
}

//Sorts the actor renderers by theie calc score
func (c *Chunk) sortARs(as []*ActorRenderer) {
	sort.Slice(as, func(i, j int) bool {
		return c.calcScore(as[i].X, as[i].Y, as[i].Z) < c.calcScore(as[j].X, as[j].Y, as[j].Z)
	})
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
	c.VisibleData = make([]Tile, 0, 100)
	//Run through all of 3d space and include
	for z := 0; z < intMin(c.MaxHeight, c.D-cutoffHeight); z++ {
		for y := c.H - 1; y >= 0; y-- { //Go backwards to go over each other
			for x := 0; x < c.W; x++ {
				tileIndex := (*c.WorldData)[z][y][x] //Minus 1 because 0 here is air not the 0 sprite
				if tileIndex != 0 {                  //again, now -1 is air
					if visible, inside := CheckVisibility(x, y, z, c.W, c.H, c.D, cutoffHeight, ChunkIndex); visible {
						//If a non-air tile is visible add it to the visible collection
						c.VisibleData = append(c.VisibleData, Tile{Sprite: pixel.NewSprite(spriteSheet, sheetFrames[tileIndex-1]), X: x, Y: y, Z: z, TileIndex: tileIndex, Hidden: inside})

					}

				}
			}
		}
	}
	//Just in case the caller didnt do this
	//Also needs to redraw everything
	c.spriteDirty = true

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

			}
			if MaxHeight != -1 {
				break
			}
		}
		if MaxHeight != -1 {
			break
		}
	}

	//Check Sprites
	for _, s := range c.SpriteDataOrdered {
		if s.Z > MaxHeight {
			MaxHeight = s.Z
		}
	}

	(*c).MaxHeight = MaxHeight + 1
}

//Add/Remove Sprite
func (c *Chunk) AddSprite(a *ActorRenderer) {

	fmt.Printf("When Added: %p\n", a)
	c.SpriteDataOrdered = append(c.SpriteDataOrdered, a)
	//Set Recalculate max, Set Re-sort Sprites

}
func (c *Chunk) RemoveSprite(a *ActorRenderer) {

	//This is really quite an ineffecient way as soon as theres more than a few sprites in the list
	if len(c.SpriteDataOrdered) > 1 {
		i := getIndex(c.SpriteDataOrdered, a)
		if i == -1 {

		}
		c.SpriteDataOrdered = removePreserveOrder(c.SpriteDataOrdered, i)
	}
	//Set Recalculate max, Set Re-sort Sprites
}

//Helpers
func removePreserveOrder(as []*ActorRenderer, s int) []*ActorRenderer {
	return append(as[:s], as[s+1:]...)
}
func getIndex(as []*ActorRenderer, value *ActorRenderer) int {
	for p, v := range as {
		if v == value {
			return p
		}
	}
	return -1
}

//Setters

//Set if it has changed
func (c *Chunk) SetChanged(changed bool) {
	c.spriteDirty = changed
}
func (c *Chunk) SetDirty(dirty bool) {
	c.tileDirty = dirty
}

//Getters
func (c *Chunk) Changed() bool {
	return c.spriteDirty
}
func (c *Chunk) Dirty() bool {
	return c.tileDirty
}
