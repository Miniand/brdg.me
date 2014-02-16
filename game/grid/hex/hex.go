package hex

import (
	"github.com/Miniand/brdg.me/game/grid"
)

const (
	DIR_NORTH = iota
	DIR_NORTH_EAST
	DIR_SOUTH_EAST
	DIR_SOUTH
	DIR_SOUTH_WEST
	DIR_NORTH_WEST
)

var neighbourOffsets = [2][6]grid.Loc{
	// Even X
	[6]grid.Loc{
		grid.Loc{0, -1},
		grid.Loc{1, -1},
		grid.Loc{1, 0},
		grid.Loc{0, 1},
		grid.Loc{-1, 0},
		grid.Loc{-1, -1},
	},
	// Odd X
	[6]grid.Loc{
		grid.Loc{0, -1},
		grid.Loc{1, 0},
		grid.Loc{1, 1},
		grid.Loc{0, 1},
		grid.Loc{-1, 1},
		grid.Loc{-1, 0},
	},
}

type Grid map[int]map[int]interface{}

func (g Grid) SetTile(l grid.Loc, tile interface{}) {
	if g[l.Y] == nil {
		g[l.Y] = map[int]interface{}{}
	}
	g[l.Y][l.X] = tile
}

func (g Grid) Tile(l grid.Loc) interface{} {
	if row := g[l.Y]; row != nil {
		return row[l.X]
	}
	return nil
}

func (g Grid) Find(tile interface{}) (grid.Loc, bool) {
	return grid.Loc{}, false
}

func (g Grid) Neighbours(l grid.Loc) []grid.Loc {
	dirs := Dirs()
	neighbours := make([]grid.Loc, len(dirs))
	for i, dir := range dirs {
		neighbours[i] = g.Neighbour(l, dir)
	}
	return neighbours
}

func (g Grid) Neighbour(l grid.Loc, dir int) grid.Loc {
	return l.Add(neighbourOffsets[l.X&1][dir])
}

func (g Grid) Bounds() (lower, upper grid.Loc) {
	for y, row := range g {
		for x, cell := range row {
			if cell != nil {
				if x < lower.X {
					lower.X = x
				}
				if x > upper.X {
					upper.X = x
				}
				if y < lower.Y {
					lower.Y = y
				}
				if y > upper.Y {
					upper.Y = y
				}
			}
		}
	}
	return
}

func Dirs() []int {
	dirs := make([]int, DIR_NORTH_WEST+1)
	for i := DIR_NORTH; i <= DIR_NORTH_WEST; i++ {
		dirs[i] = i
	}
	return dirs
}
