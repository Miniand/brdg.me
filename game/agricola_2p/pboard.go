package agricola_2p

import "github.com/Miniand/brdg.me/game/cost"

type PBoard struct {
	Resources    cost.Cost
	Tiles        Tiles
	XStart, XEnd int
}

func NewPBoard() *PBoard {
	return &PBoard{
		Resources: cost.Cost{
			ResourceBorder: 9,
		},
		Tiles: Tiles{
			Loc{0, 2}: {
				Building: Cottage{},
			},
			// Testing
			Loc{0, 0}: {
				Borders: Up | Left,
			},
			Loc{0, 1}: {
				Borders: Left,
			},
			Loc{1, 1}: {
				Borders: Up,
			},
			Loc{2, 1}: {
				Borders: Left,
			},
			Loc{1, 2}: {
				Borders: Up,
			},
			Loc{1, 0}: {
				Borders: Left,
			},
		},
		XEnd: 1,
	}
}
