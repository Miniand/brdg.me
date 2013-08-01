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
	dirStopped := map[int]map[int]bool{}
	dirStopped[-1] = map[int]bool{}
	dirStopped[1] = map[int]bool{}
	for dist := 1; dist <= 8; dist++ {
		for _, rankDir := range []int{-1, 1} {
			for _, fileDir := range []int{-1, 1} {
				l := Location{
					from.File + dist*fileDir,
					from.Rank + dist*rankDir}
				if !dirStopped[rankDir][fileDir] && IsValidLocation(l) {
					piece := b.PieceAt(l)
					if piece != nil {
						dirStopped[rankDir][fileDir] = true
						if piece.GetTeam() != bi.Team {
							// Can take this piece
							to = append(to, l)
						}
					} else {
						// Empty space
						to = append(to, l)
					}
				}
			}
		}
	}
	return
}
