package alhambra

import "github.com/Miniand/brdg.me/render"

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

var TileAbbrLen = len(TileAbbrs[TileTypeEmpty])

type Tile struct {
	Type  int
	Cost  int
	Walls map[int]bool
}
