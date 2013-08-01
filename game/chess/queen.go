package chess

type Queen struct {
	Piece
}

func (q Queen) Rune() rune {
	if q.Team == TEAM_WHITE {
		return '♕'
	}
	return '♛'
}

func (q Queen) AvailableMoves(from Location, b Board) (to []Location) {
	return
}
