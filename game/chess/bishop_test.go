package chess

import (
	"testing"
)

func TestDiagonalMovement(t *testing.T) {
	b := parseBoard(`
········
········
····♛···
········
··♗·····
········
········
·····♘··`, t)
	l := Location{FILE_C, RANK_4}
	bi := b.PieceAt(l).(*Bishop)
	to := bi.AvailableMoves(l, b)
	checkLocations := []Location{
		Location{FILE_A, RANK_2},
		Location{FILE_B, RANK_3},
		Location{FILE_B, RANK_5},
		Location{FILE_A, RANK_6},
		Location{FILE_D, RANK_5},
		Location{FILE_E, RANK_6},
		Location{FILE_D, RANK_3},
		Location{FILE_E, RANK_2},
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
