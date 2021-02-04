package render


import(
	"fmt"
)
//A type to hold the chunk data in one place
type Chunk struct{
	changed bool
	MaxHeight int
	WorldData *[][][]int
	W,H,D int
}


//Calculators
//Calculate the maximum non air block
func (c *Chunk) CalculateMax(){
	//Run down from top and stop once a non air block is hit
	MaxHeight:=-1

	for z:=c.D-1; z>=0; z--{
		for y := 0; y < c.H; y++ { 
			for x := 0; x < c.W; x++ {
				if (*c.WorldData)[z][y][x]!=0{
					MaxHeight=z
					break
				} 
			}
			if MaxHeight!=-1{
				break
			}
		}
		if MaxHeight!=-1{
			break
		}
	}
	fmt.Println("Caluclated Max Height",MaxHeight)
	c.MaxHeight=MaxHeight
}


//Setters
//Set if it has changed
func (c *Chunk) SetChanged( changed bool){
	c.changed=changed
}

//Getters
func (c *Chunk) GetChanged() bool {
	return c.changed
}

