package cathedral

type PlayerType struct {
	Player, Type int
}

type Location struct {
	X, Y int
}

type Piece struct {
	PlayerType
	Positions   []Location
	Directional bool
}

var Pieces = map[PlayerType]Piece{
	// Player 1
	PlayerType{0, 0}: Piece{
		PlayerType{0, 0},
		[]Location{
			{0, 0},
			{0, 1},
			{-1, 1},
			{0, 2},
			{1, 2},
		},
		true,
	},
	PlayerType{0, 1}: Piece{
		PlayerType{0, 1},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{1, 2},
			{2, 2},
		},
		true,
	},
	PlayerType{0, 2}: Piece{
		PlayerType{0, 2},
		[]Location{
			{0, 0},
			{0, 1},
			{-1, 1},
			{1, 1},
			{0, 2},
		},
		false,
	},
	PlayerType{0, 3}: Piece{
		PlayerType{0, 3},
		[]Location{
			{0, 0},
			{1, 0},
			{0, 1},
			{0, 2},
			{1, 2},
		},
		true,
	},
	PlayerType{0, 4}: Piece{
		PlayerType{0, 4},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{0, 2},
		},
		true,
	},
	PlayerType{0, 5}: Piece{
		PlayerType{0, 5},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{1, 2},
		},
		true,
	},
	PlayerType{0, 6}: Piece{
		PlayerType{0, 6},
		[]Location{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 1},
		},
		false,
	},
	PlayerType{0, 7}: Piece{
		PlayerType{0, 7},
		[]Location{
			{0, 0},
			{0, 1},
			{0, 2},
		},
		true,
	},
	PlayerType{0, 8}: Piece{
		PlayerType{0, 8},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
		},
		true,
	},
	PlayerType{0, 9}: Piece{
		PlayerType{0, 9},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
		},
		true,
	},
	PlayerType{0, 10}: Piece{
		PlayerType{0, 10},
		[]Location{
			{0, 0},
			{0, 1},
		},
		true,
	},
	PlayerType{0, 11}: Piece{
		PlayerType{0, 11},
		[]Location{
			{0, 0},
			{0, 1},
		},
		true,
	},
	PlayerType{0, 12}: Piece{
		PlayerType{0, 12},
		[]Location{
			{0, 0},
		},
		false,
	},
	PlayerType{0, 13}: Piece{
		PlayerType{0, 13},
		[]Location{
			{0, 0},
		},
		false,
	},

	// Player 2
	PlayerType{1, 0}: Piece{
		PlayerType{1, 0},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{0, 2},
			{-1, 2},
		},
		true,
	},
	PlayerType{1, 1}: Piece{
		PlayerType{1, 1},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{1, 2},
			{2, 2},
		},
		true,
	},
	PlayerType{1, 2}: Piece{
		PlayerType{1, 2},
		[]Location{
			{0, 0},
			{0, 1},
			{-1, 1},
			{1, 1},
			{0, 2},
		},
		false,
	},
	PlayerType{1, 3}: Piece{
		PlayerType{1, 3},
		[]Location{
			{0, 0},
			{1, 0},
			{0, 1},
			{0, 2},
			{1, 2},
		},
		true,
	},
	PlayerType{1, 4}: Piece{
		PlayerType{1, 4},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
			{0, 2},
		},
		true,
	},
	PlayerType{1, 5}: Piece{
		PlayerType{1, 5},
		[]Location{
			{0, 0},
			{0, 1},
			{-1, 1},
			{-1, 2},
		},
		true,
	},
	PlayerType{1, 6}: Piece{
		PlayerType{1, 6},
		[]Location{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 1},
		},
		false,
	},
	PlayerType{1, 7}: Piece{
		PlayerType{1, 7},
		[]Location{
			{0, 0},
			{0, 1},
			{0, 2},
		},
		true,
	},
	PlayerType{1, 8}: Piece{
		PlayerType{1, 8},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
		},
		true,
	},
	PlayerType{1, 9}: Piece{
		PlayerType{1, 9},
		[]Location{
			{0, 0},
			{0, 1},
			{1, 1},
		},
		true,
	},
	PlayerType{1, 10}: Piece{
		PlayerType{1, 10},
		[]Location{
			{0, 0},
			{0, 1},
		},
		true,
	},
	PlayerType{1, 11}: Piece{
		PlayerType{1, 11},
		[]Location{
			{0, 0},
			{0, 1},
		},
		true,
	},
	PlayerType{1, 12}: Piece{
		PlayerType{1, 12},
		[]Location{
			{0, 0},
		},
		false,
	},
	PlayerType{1, 13}: Piece{
		PlayerType{1, 13},
		[]Location{
			{0, 0},
		},
		false,
	},

	// Cathedral
	PlayerType{2, 0}: Piece{
		PlayerType{2, 0},
		[]Location{
			{0, 0},
			{0, 1},
			{0, 2},
			{-1, 2},
			{1, 2},
			{0, 3},
		},
		true,
	},
}
