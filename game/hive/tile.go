package hive

const (
	TILE_QUEEN_BEE = iota
	TILE_BEETLE
	TILE_GRASSHOPPER
	TILE_SPIDER
	TILE_SOLDIER_ANT
)

var tileNames = map[int]string{
	TILE_QUEEN_BEE:   "queen",
	TILE_BEETLE:      "beetle",
	TILE_GRASSHOPPER: "grasshopper",
	TILE_SPIDER:      "spider",
	TILE_SOLDIER_ANT: "soldier ant",
}

var tileShortNames = map[int]string{
	TILE_QUEEN_BEE:   "queen",
	TILE_BEETLE:      "beetl",
	TILE_GRASSHOPPER: "hoppr",
	TILE_SPIDER:      "spdr",
	TILE_SOLDIER_ANT: "ant",
}

type Tile struct {
	Type   int
	Player int
}

func (t *Tile) Colour() string {
	if t.Player == 0 {
		return "blue"
	}
	return "red"
}

func (t *Tile) Message() string {
	return tileShortNames[t.Type]
}
