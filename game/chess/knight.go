package chess

type Knight struct {
	Piece
}

func (k Knight) Rune() rune {
	if k.Team == TEAM_WHITE {
		return '♘'
	}
	return '♞'
}

func (k Knight) AvailableMoves(from Location, b Board) (to []Location) {
	return
}
