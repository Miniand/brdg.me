package chess

import (
	"testing"
)

func TestQueenMovement(t *testing.T) {
	b := parseBoard(`
········
········
··♛·····
···♘····
··♕···♘·
········
♛·······
········`, t)
	l := Location{FILE_C, RANK_4}
	q := b.PieceAt(l).(*Queen)
	to := q.AvailableMoves(l, b)
	checkLocations := []Location{
		// Rook moves
		Location{FILE_C, RANK_1},
		Location{FILE_C, RANK_2},
		Location{FILE_C, RANK_3},
		Location{FILE_C, RANK_5},
		Location{FILE_C, RANK_6},
		Location{FILE_A, RANK_4},
		Location{FILE_B, RANK_4},
		Location{FILE_D, RANK_4},
		Location{FILE_E, RANK_4},
		Location{FILE_F, RANK_4},
		// Bishop moves
		Location{FILE_B, RANK_3},
		Location{FILE_A, RANK_2},
		Location{FILE_B, RANK_5},
		Location{FILE_A, RANK_6},
		Location{FILE_D, RANK_3},
		Location{FILE_E, RANK_2},
		Location{FILE_F, RANK_1},
	}
	if len(to) != len(checkLocations) {
		t.Fatalf("Number of available moves doesn't match expected\nAvailable:\n%#v\nExpected:\n%#v", to, checkLocations)
	}
	for _, check := range checkLocations {
		if !inLocations(check, to) {
			t.Fatalf("Available:\n%#v\nExpected:\n%#v", to, checkLocations)
		}
	}
}
