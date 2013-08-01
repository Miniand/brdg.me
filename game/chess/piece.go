package chess

type Piecer interface {
	Rune() rune
	AvailableMoves(from Location, b Board) (to []Location)
	GetTeam() int
}

type Piece struct {
	Team int
}

func (p Piece) GetTeam() int {
	return p.Team
}
