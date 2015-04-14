package cathedral

import (
	"strconv"

	"github.com/Miniand/brdg.me/render"
)

type PlayerType struct {
	Player, Type int
}

type Piece struct {
	PlayerType
	Positions   Locs
	Directional bool
}

func (p Piece) TileAt(loc Loc) (Tile, bool) {
	origin := true
	for _, l := range p.Positions {
		if l == loc {
			t := Tile{
				PlayerType: p.PlayerType,
			}
			if origin {
				t.Text = strconv.Itoa(p.Type)
			}
			return t, true
		}
		origin = false
	}
	return Tile{}, false
}

func (p Piece) Bounds() (lower, upper Loc) {
	first := true
	for _, l := range p.Positions {
		if first || l.X < lower.X {
			lower.X = l.X
		}
		if first || l.Y < lower.Y {
			lower.Y = l.Y
		}
		if first || l.X > upper.X {
			upper.X = l.X
		}
		if first || l.Y > upper.Y {
			upper.Y = l.Y
		}
		first = false
	}
	return
}

func (p Piece) Width() int {
	lower, upper := p.Bounds()
	return upper.Sub(lower).X + 1
}

func (p Piece) Render() string {
	cells := [][]interface{}{}
	lower, upper := p.Bounds()
	for y := lower.Y; y <= upper.Y; y++ {
		row := []interface{}{}
		for x := lower.X; x <= upper.X; x++ {
			rt, _ := RenderTile(p, Loc{x, y})
			row = append(row, rt)
		}
		cells = append(cells, row)
	}
	return render.Table(cells, 0, 0)
}

var Pieces = map[int][]Piece{
	// Player 1
	0: {
		Piece{
			PlayerType{0, 1},
			Locs{
				{0, 0},
				{0, 1},
				{-1, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 2},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
				{2, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 3},
			Locs{
				{0, 0},
				{0, 1},
				{-1, 1},
				{1, 1},
				{0, 2},
			},
			false,
		},
		Piece{
			PlayerType{0, 4},
			Locs{
				{0, 0},
				{1, 0},
				{0, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 5},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 6},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 7},
			Locs{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
			false,
		},
		Piece{
			PlayerType{0, 8},
			Locs{
				{0, 0},
				{0, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 9},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 10},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 11},
			Locs{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 12},
			Locs{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 13},
			Locs{
				{0, 0},
			},
			false,
		},
		Piece{
			PlayerType{0, 14},
			Locs{
				{0, 0},
			},
			false,
		},
	},

	// Player 2
	1: {
		// Cathedral first
		Piece{
			PlayerType{PlayerCathedral, 1},
			Locs{
				{0, 0},
				{0, 1},
				{0, 2},
				{-1, 2},
				{1, 2},
				{0, 3},
			},
			true,
		},
		Piece{
			PlayerType{1, 2},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
				{-1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 3},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
				{2, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 4},
			Locs{
				{0, 0},
				{0, 1},
				{-1, 1},
				{1, 1},
				{0, 2},
			},
			false,
		},
		Piece{
			PlayerType{1, 5},
			Locs{
				{0, 0},
				{1, 0},
				{0, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 6},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 7},
			Locs{
				{0, 0},
				{0, 1},
				{-1, 1},
				{-1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 8},
			Locs{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
			false,
		},
		Piece{
			PlayerType{1, 9},
			Locs{
				{0, 0},
				{0, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 10},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 11},
			Locs{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 12},
			Locs{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 13},
			Locs{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 14},
			Locs{
				{0, 0},
			},
			false,
		},
		Piece{
			PlayerType{1, 15},
			Locs{
				{0, 0},
			},
			false,
		},
	},
}
