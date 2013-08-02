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
	for _, dir := range [][]int{
		[]int{1, 1},
		[]int{1, -1},
		[]int{-1, 1},
		[]int{-1, -1},
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
				if piece.GetTeam() != bi.Team {
					to = append(to, l)
				}
				break
			}
			to = append(to, l)
		}
	}
	return
}
