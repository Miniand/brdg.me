package cathedral

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
