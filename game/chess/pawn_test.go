package chess

import (
	"testing"
)

func TestPawnIsPiecer(t *testing.T) {
	var p Piecer
	p = Pawn{}
	if p == nil {
		t.Fatal("r is nil")
	}
}

func TestPawnAdvanceOne(t *testing.T) {
	b := Board{}
	p := Pawn{
		Piece{
			Team: TEAM_WHITE,
		},
	}
	to := p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if !inLocations(Location{FILE_B, RANK_5}, to) {
		t.Fatal("Did not find move to advance one")
	}
	b.Squares[FILE_B][RANK_5] = &Pawn{}
	to = p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if inLocations(Location{FILE_B, RANK_5}, to) {
		t.Fatal("Found move to advance one when there was a piece there")
	}
}
