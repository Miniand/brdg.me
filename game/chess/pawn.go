package chess

type Pawn struct {
	Piece
}

func (p Pawn) Rune() rune {
	if p.Team == TEAM_WHITE {
		return '♙'
	}
	return '♟'
}

func (p Pawn) AvailableMoves(from Location, b Board) (to []Location) {
	advanceOne := Location{from.File, from.Rank + p.Team}
	if IsValidLocation(advanceOne) && b.IsEmpty(advanceOne) {
		to = append(to, advanceOne)
		if from.Rank == StartRank(p.Team)+p.Team {
			advanceTwo := Location{advanceOne.File, advanceOne.Rank + p.Team}
			if b.IsEmpty(advanceTwo) {
				to = append(to, advanceTwo)
			}
		}
	}
	return
}
