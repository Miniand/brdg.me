package acquire

import (
	"github.com/Miniand/brdg.me/game/card"
)

type Tile struct {
	Row, Column int
}

// Sort by column first, then row
func (t Tile) Compare(other card.Comparer) (int, bool) {
	otherT, ok := other.(Tile)
	if !ok {
		// Different types
		return 0, false
	}
	if t.Column < otherT.Column {
		return -1, true
	} else if t.Column > otherT.Column {
		return 1, true
	} else if t.Row < otherT.Row {
		return -1, true
	} else if t.Row > otherT.Row {
		return 1, true
	}
	return 0, true
}
