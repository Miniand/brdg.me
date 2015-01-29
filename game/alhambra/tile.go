package alhambra

import "github.com/Miniand/brdg.me/render"

const (
	TileTypeEmpty = iota
	TileTypePavillion
	TileTypeSeraglio
	TileTypeArcades
	TileTypeChambers
	TileTypeGarden
	TileTypeTower
)

var TileAbbrs = map[int]string{
	TileTypeEmpty:     "  ",
	TileTypePavillion: "Pa",
	TileTypeSeraglio:  "Se",
	TileTypeArcades:   "Ar",
	TileTypeChambers:  "Ch",
	TileTypeGarden:    "Ga",
	TileTypeTower:     "To",
}

var TileColours = map[int]string{
	TileTypeEmpty:     render.Black,
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
