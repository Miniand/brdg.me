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
	// Queen is just a mix of rook and bishop moves
	r := Rook{}
	r.Team = q.Team
	to = r.AvailableMoves(from, b)
	bi := Bishop{}
	bi.Team = q.Team
	to = append(to, bi.AvailableMoves(from, b)...)
	return
}
