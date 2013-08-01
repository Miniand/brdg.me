package chess

type Rook struct {
	Piece
}

func (r Rook) Rune() rune {
	if r.Team == TEAM_WHITE {
		return '♖'
	}
	return '♜'
}

func (r Rook) AvailableMoves(from Location, b Board) (to []Location) {
	return
}
