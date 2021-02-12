package main

import(
		"math"
		"math/rand"
		
		"./Config"
		"./Materials"
)

//Returns a value for the surface based on a mathematical function of the form f(x,y) defined here
func landFunc(x, y int) int {
	w := 32
	h := 40
	r := math.Sqrt(float64((x-w)*(x-w))+float64((y-h)*(y-h))) / 30
	t := math.Atan2(float64(x-w), float64(y-h))
	f := math.Cos(t * 6)

	v := 0.0
	if f-r < 4 {
		v = f - r
	}
	
	v*=1+r
	return 6 + int(v)/2

}
/*
//GenMap2 generates a 3d world of the form z=f(x,y) approx where f is landFunc() and tiles selected based off of z
func GenMap2(w, h, d int) [][][]int {

	world := make([][][]int, 0, d)

	for z := 0; z < d; z++ {
		floor := make([][]int, 0, h)

		for y := 0; y < h; y++ {
			row := make([]int, 0, w)
			for x := 0; x < w; x++ {
				block := materials.STONE_1 + rand.Intn(2) 
				if z > landFunc(x, y) {
					block = materials.AIR
				} else if z == landFunc(x, (h-1)-y) {
					block = materials.GRASS + rand.Intn(8)
				}

				row = append(row, block)

			}
			floor = append(floor, row)

		}
		world = append(world, floor)
	}
	//fmt.Println(world)
	return world
}
*/

//GenMap3 creates a map  in the same way as GenMap2 but that accounts for chunks
func GenMap3(w, h, d, chunkx, chunky int) [config.ChunkDepth][config.ChunkHeight][config.ChunkWidth]int {

	world := [config.ChunkDepth][config.ChunkHeight][config.ChunkWidth]int{}

	for z := 0; z < d; z++ {
		floor:=[config.ChunkHeight][config.ChunkWidth]int{}
		for y := 0; y < h; y++ {
			row := [config.ChunkWidth]int{}//make([]int, 0, w)
			for x := 0; x < w; x++ {
				val := landFunc(x+chunkx*16, y+chunky*16)
				block := materials.Sprites["STONE_1"] + rand.Intn(2) //Default Block
				if z > val {
					block = materials.Sprites["AIR"]
				} 
				if z == val {
					block = materials.Sprites["GRASS"]+ rand.Intn(8)
				}
				row[x]=block


			}
			floor[y]=row

		}
		world [z]=floor
	}
	return world
}
