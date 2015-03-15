package cathedral

import "github.com/Miniand/brdg.me/render"

type PlayerType struct {
	Player, Type int
}

type Location struct {
	X, Y int
}

func (l Location) Rotate(n int) Location {
	switch {
	case n > 0:
		return (Location{-l.Y, l.X}).Rotate(n - 1)
	case n < 0:
		return (Location{l.Y, -l.X}).Rotate(n + 1)
	default:
		return l
	}
}

type Locations []Location

func (ls Locations) Rotate(n int) Locations {
	nls := make(Locations, len(ls))
	for i, l := range ls {
		nls[i] = l.Rotate(n)
	}
	return nls
}

type Piece struct {
	PlayerType
	Positions   Locations
	Directional bool
}

func (p Piece) TileAt(x, y int) (Tile, bool) {
	for _, l := range p.Positions {
		if x == l.X && y == l.Y {
			return Tile{
				PlayerType: p.PlayerType,
			}, true
		}
	}
	return Tile{}, false
}

func (p Piece) Bounds() (x1, y1, x2, y2 int) {
	first := true
	for _, l := range p.Positions {
		if first || l.X < x1 {
			x1 = l.X
		}
		if first || l.Y < y1 {
			y1 = l.Y
		}
		if first || l.X > x2 {
			x2 = l.X
		}
		if first || l.Y > y2 {
			y2 = l.Y
		}
		first = false
	}
	return
}

func (p Piece) Width() int {
	x1, _, x2, _ := p.Bounds()
	return x2 - x1 + 1
}

func (p Piece) Render() string {
	cells := [][]interface{}{}
	xMin, yMin, xMax, yMax := p.Bounds()
	for y := yMin; y <= yMax; y++ {
		row := []interface{}{}
		for x := xMin; x <= xMax; x++ {
			rt, _ := RenderTile(p, x, y)
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
			PlayerType{0, 0},
			Locations{
				{0, 0},
				{0, 1},
				{-1, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 1},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
				{2, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 2},
			Locations{
				{0, 0},
				{0, 1},
				{-1, 1},
				{1, 1},
				{0, 2},
			},
			false,
		},
		Piece{
			PlayerType{0, 3},
			Locations{
				{0, 0},
				{1, 0},
				{0, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 4},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 5},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 6},
			Locations{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
			false,
		},
		Piece{
			PlayerType{0, 7},
			Locations{
				{0, 0},
				{0, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{0, 8},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 9},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 10},
			Locations{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 11},
			Locations{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{0, 12},
			Locations{
				{0, 0},
			},
			false,
		},
		Piece{
			PlayerType{0, 13},
			Locations{
				{0, 0},
			},
			false,
		},
	},

	// Player 2
	1: {
		// Cathedral first
		Piece{
			PlayerType{PlayerCathedral, 0},
			Locations{
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
			PlayerType{1, 0},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
				{-1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 1},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 2},
				{2, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 2},
			Locations{
				{0, 0},
				{0, 1},
				{-1, 1},
				{1, 1},
				{0, 2},
			},
			false,
		},
		Piece{
			PlayerType{1, 3},
			Locations{
				{0, 0},
				{1, 0},
				{0, 1},
				{0, 2},
				{1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 4},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 5},
			Locations{
				{0, 0},
				{0, 1},
				{-1, 1},
				{-1, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 6},
			Locations{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
			false,
		},
		Piece{
			PlayerType{1, 7},
			Locations{
				{0, 0},
				{0, 1},
				{0, 2},
			},
			true,
		},
		Piece{
			PlayerType{1, 8},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 9},
			Locations{
				{0, 0},
				{0, 1},
				{1, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 10},
			Locations{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 11},
			Locations{
				{0, 0},
				{0, 1},
			},
			true,
		},
		Piece{
			PlayerType{1, 12},
			Locations{
				{0, 0},
			},
			false,
		},
		Piece{
			PlayerType{1, 13},
			Locations{
				{0, 0},
			},
			false,
		},
	},
}
