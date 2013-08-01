package chess

type Bishop struct {
	Piece
}

func (bi Bishop) Rune() rune {
	if bi.Team == TEAM_WHITE {
		return '♗'
	}
	return '♝'
}

func (bi Bishop) AvailableMoves(from Location, b Board) (to []Location) {
	return
}
