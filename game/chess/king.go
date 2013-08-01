package chess

type King struct {
	Piece
}

func (k King) Rune() rune {
	if k.Team == TEAM_WHITE {
		return '♔'
	}
	return '♚'
}

func (k King) AvailableMoves(from Location, b Board) (to []Location) {
	return
}
