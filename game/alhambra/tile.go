package alhambra

const (
	TileTypeEmpty = iota
	TileTypeBlah
)

var TileRunes = map[int]rune{
	TileTypeEmpty: ' ',
	TileTypeBlah:  'B',
}

type Tile struct {
	Type  int
	Cost  int
	Walls map[int]bool
}
