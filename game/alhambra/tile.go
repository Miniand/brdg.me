package alhambra

import (
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/render"
)

const (
	TileTypeEmpty = iota
	TileTypeFountain
	TileTypePavillion
	TileTypeSeraglio
	TileTypeArcades
	TileTypeChambers
	TileTypeGarden
	TileTypeTower
)

var TileAbbrs = map[int]string{
	TileTypeEmpty:     "   ",
	TileTypeFountain:  " F ",
	TileTypePavillion: "Pav",
	TileTypeSeraglio:  "Ser",
	TileTypeArcades:   "Arc",
	TileTypeChambers:  "Cha",
	TileTypeGarden:    "Gar",
	TileTypeTower:     "Tow",
}

var TileColours = map[int]string{
	TileTypeEmpty:     render.Black,
	TileTypeFountain:  render.Gray,
	TileTypePavillion: render.Cyan,
	TileTypeSeraglio:  render.Red,
	TileTypeArcades:   render.Blue,
	TileTypeChambers:  render.Yellow,
	TileTypeGarden:    render.Green,
	TileTypeTower:     render.Magenta,
}

var ScoringTileTypes = []int{
	TileTypePavillion,
	TileTypeSeraglio,
	TileTypeArcades,
	TileTypeChambers,
	TileTypeGarden,
	TileTypeTower,
}

var TileAbbrLen = len(TileAbbrs[TileTypeEmpty])
var TileWidth = TileAbbrLen + 1

type Tile struct {
	Type  int
	Cost  int
	Walls map[int]bool
}

func ShuffleTiles(tiles []Tile) []Tile {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(tiles)
	shuffled := make([]Tile, l)
	for i, k := range r.Perm(l) {
		shuffled[i] = tiles[k]
	}
	return shuffled
}

func NotEmpty(tiles []Tile) []Tile {
	notEmpty := []Tile{}
	for _, t := range tiles {
		if t.Type != TileTypeEmpty {
			notEmpty = append(notEmpty, t)
		}
	}
	return notEmpty
}
