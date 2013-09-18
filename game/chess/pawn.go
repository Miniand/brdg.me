package chess

type Pawn struct {
	Piece
	MoveWasAdvanceTwo bool // Used for en passant check, @see https://en.wikipedia.org/wiki/En_passant
}

func (p Pawn) Rune() rune {
	if p.Team == TEAM_WHITE {
		return '♙'
	}
	return '♟'
}

func (p Pawn) AvailableMoves(from Location, b Board) (to []Move) {
	// Check advancing one space forward
	advanceOne := Location{from.File, from.Rank + p.Team}
	if IsValidLocation(advanceOne) && b.IsEmpty(advanceOne) {
		to = append(to, Move{
			From: from,
			To:   advanceOne,
		})
		// Check if it's the initial move, which allows two spaces
		if from.Rank == StartRank(p.Team)+p.Team {
			advanceTwo := Location{advanceOne.File, advanceOne.Rank + p.Team}
			if b.IsEmpty(advanceTwo) {
				to = append(to, Move{
					From: from,
					To:   advanceTwo,
				})
			}
		}
	}
	// Check if there is a piece diagonally forward
	for _, l := range []Location{
		Location{from.File - 1, from.Rank + p.Team},
		Location{from.File + 1, from.Rank + p.Team},
	} {
		if IsValidLocation(l) {
			attackDownPiece := b.PieceAt(l)
			if attackDownPiece != nil {
				if attackDownPiece.GetTeam() != p.Team {
					to = append(to, Move{
						From:   from,
						To:     l,
						TakeAt: &l,
					})
				}
			} else {
				// Check en passant, @see https://en.wikipedia.org/wiki/En_passant
				adjacentPiece := b.PieceAt(Location{l.File, from.Rank})
				adjacentPawn, ok := adjacentPiece.(*Pawn)
				if ok && adjacentPawn.Team != p.Team &&
					adjacentPawn.MoveWasAdvanceTwo {
					to = append(to, Move{
						From:   from,
						To:     l,
						TakeAt: &Location{l.File, from.Rank},
					})
				}
			}
		}
	}
	return
}
