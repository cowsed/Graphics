package materials

import (
	"encoding/csv"
	_"fmt"
	"github.com/faiface/pixel"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var NumTiles = 0
var Sprites map[string]int
var SpritesByIndex []string
var Descriptions []string


var Picture pixel.Picture
var frames []pixel.Rect

const AIR = 0

func GetData() (pixel.Picture, []pixel.Rect) {
	return Picture, frames
}

func LoadSprites(fname string, picFname string) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(strings.NewReader(string(content)))

	frames = make([]pixel.Rect, 0, NumTiles)
	Sprites = make(map[string]int)
	Descriptions=make([]string, 0, NumTiles)

	index := 1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(record)

		Sprites[record[0]] = index
		Descriptions = append(Descriptions, record[9])

		SpritesByIndex=append(SpritesByIndex,record[0])

		x, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			panic("Theres probably a space in your numbers")
		}

		y, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			panic("Theres probably a space in your numbers")
		}
		w, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			panic("Theres probably a space in your numbers")
		}
		h, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			panic("Theres probably a space in your numbers")
		}

		frames = append(frames, pixel.R(x, y, x+w, y+h))

		index++
	}

	Picture, err = loadPicture(picFname)
	if err != nil {
		log.Fatal(err)
	}

}

//Loads the picture
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

//Material tile offsetts to reference

/*
const (
	none= iota
	CURSOR_BLOCK
	CURSOR_TOP
	GRASS
	GRASS_LESS_1
	GRASS_LESS_2

	GRASS_NONE

	GRASS_FLOWER_1
	GRASS_FLOWER_2
	GRASS_FLOWER_3
	GRASS_FLOWER_4

	STONE_1
	STONE_2

	BRICK
	TREE_TRUNK

	ROCK_WALL_V_1
	ROCK_WALL_H_1
	ROCK_WALL_FR
)
*/
