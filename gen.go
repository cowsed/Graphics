package main

import "fmt"

import "math"
import "math/rand"
import "./Materials"

func MakeFlags(w int, h int, x int, y int, m [][]int) string {
	flag := []bool{true, true, true, true, true, true, true, true, true} //representative of if its safe to check
	if y == 0 {
		flag[1] = false
		flag[2] = false
		flag[8] = false
	} else if y == h-1 {
		flag[4] = false
		flag[5] = false
		flag[6] = false
	}
	if x == 0 {
		flag[6] = false
		flag[7] = false
		flag[8] = false
	} else if x == w-1 {
		flag[2] = false
		flag[3] = false
		flag[4] = false
	}

	cm := [...][2]int{[2]int{0, 0}, [2]int{0, -1}, [2]int{1, -1}, [2]int{1, 0}, [2]int{1, 1}, [2]int{0, 1}, [2]int{-1, 1}, [2]int{-1, 0}, [2]int{-1, -1}}

	s := ""

	for i, p := range cm {
		if flag[i] {
			r := m[y+p[1]][x+p[0]]
			if r == 1 {
				s += "1"
			} else {
				s += "0"
			}

		} else {
			s += "1"
		}
	}
	fmt.Println("S: ", s)

	return s
}

/*

func makeArrangement(flags string) int {
	//flags  starting from back going clockwise 10101010 is 1s at edges nothing at corners
	//assume ends are 1s
	if flags[0]=='1'{return HILL_C}
	switch flags{
		case ""://,"110101010","110001000","100100010"://Surrounded on at least 2 opposite sides
			return HILL_C
		case "10000000":
			return HILL_F
		case "100101000":
			return HILL_FR2
		case "110000010":
			return HILL_R
		case "110000000":
			return HILL_FR
		case "011101111":
			return HILL_BR
		case "111110101":
			return HILL_FL

		case "101111101":
			return HILL_BL2
		case "101011111":
			return HILL_BR2

		default:
			return GRASS_1
	}
}

func GenHillSmart() [][][]int {

	noise:=[][]int{
				[]int{1,0,1,0},
				[]int{1,1,0,0},
				[]int{0,1,1,0},
				}

	h:=len(noise)
	w:=len(noise[0])

	floor := make([][]int, 0)
	for y := 0; y < h; y++ {
		row := make([]int, 0)
		for x := 0; x < w; x++ {
			//Flag making
			fs:=MakeFlags(w,h,x,y,noise)
			//fmt.Println(fs)

			row = append(row, makeArrangement(fs))
		}
		floor = append(floor, row)

	}
	return [][][]int{floor}

}
*/

func GenMap1(w, h, d int) [][][]int {
	world := [][][]int{}
	i := 0
	for z := 0; z < d; z++ {
		floor := make([][]int, 0)

		for y := 0; y < h; y++ {
			row := make([]int, 0)
			for x := 0; x < w; x++ {
				row = append(row, materials.HILL_C)
				i++
				i = i % 130
				fmt.Println(i)
			}
			floor = append(floor, row)

		}
		world = append(world, floor)
	}
	//fmt.Println(world)
	return world
}

/*
func GenMap4() [][][]int {
	floor:=[][]int{}
	for y := 0; y < 12; y++ {
		floor = append(floor, []int{HILL_BL2, HILL_B, HILL_BR2,HILL_BL, HILL_B, HILL_BR,HILL_BL, HILL_B, HILL_BR})
		floor = append(floor, []int{HILL_L, HILL_C, HILL_R,HILL_L, HILL_C, HILL_R,HILL_L, HILL_C, HILL_R})
		floor = append(floor, []int{HILL_FL2, HILL_F, HILL_FR2,HILL_FL, HILL_F, HILL_FR,HILL_FL, HILL_F, HILL_FR})

	}

	return [][][]int{floor}

}

func GenMap3() [][][]int {

	return [][][]int{
		[][]int{
			[]int{HILL_C, HILL_C, HILL_C},
			[]int{HILL_C, 0, HILL_C},
			[]int{HILL_C, HILL_C, HILL_C},
		},
		[][]int{
			[]int{HILL_BL, HILL_B, HILL_BR},
			[]int{HILL_L, HILL_C, HILL_R},
			[]int{HILL_FL, HILL_F, HILL_FR},
		},
	}
}
*/

func GenMap2(w, h, d int) [][][]int {
	freq := .5

	world := [][][]int{}

	for z := 0; z < d; z++ {
		floor := make([][]int, 0)

		for y := 0; y < h; y++ {
			row := make([]int, 0)
			for x := 0; x < w; x++ {
				block := 160
				if z > int(3+2*math.Sin(float64(x)*freq)*math.Sin(float64(y)*freq)) {
					block = 0
				} else if z == int(3+2*math.Sin(float64(x)*freq)*math.Sin(float64(y)*freq)) {
					block = materials.HILL_C
				} else if z == 0 {
					block = materials.ROCK + rand.Intn(2)
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
