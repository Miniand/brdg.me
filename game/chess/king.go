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

func (k King) IsInCheck(at Location, b Board) bool {
	for file := FILE_A; file <= FILE_H; file++ {
		for rank := RANK_1; rank <= RANK_8; rank++ {
			loc := Location{file, rank}
			if piece := b.PieceAt(loc); piece != nil &&
				piece.GetTeam() != k.GetTeam() {
				for _, move := range piece.AvailableMoves(loc, b) {
					if *move.TakeAt == at {
						return true
					}
				}
			}
		}
	}
	return false
}

func (k King) AvailableMoves(from Location, b Board) (to []Move) {
	return
}
