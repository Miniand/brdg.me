package alhambra

const (
	TileTypeEmpty = iota
	TileTypeBlah
)

var TileStrs = map[int]string{
	TileTypeEmpty: " ",
	TileTypeBlah:  "B",
}

type Tile struct {
	Type  int
	Cost  int
	Walls map[int]bool
}
