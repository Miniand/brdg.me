package chess

type Rook struct {
	Piece
	HasMoved bool
}

func (r Rook) Rune() rune {
	if r.Team == TEAM_WHITE {
		return '♖'
	}
	return '♜'
}

func (r Rook) AvailableMoves(from Location, b Board) (to []Location) {
	for _, dir := range [][]int{
		[]int{1, 0},
		[]int{-1, 0},
		[]int{0, 1},
		[]int{0, -1},
	} {
		for dist := 1; dist < 8; dist++ {
			l := Location{
				from.File + dist*dir[0],
				from.Rank + dist*dir[1]}
			if !IsValidLocation(l) {
				break
			}
			piece := b.PieceAt(l)
			if piece != nil {
				if piece.GetTeam() != r.Team {
					to = append(to, l)
				}
				break
			}
			to = append(to, l)
		}
	}
	return
}
