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
	p := Pawn{}
	p.Team = TEAM_WHITE
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

func TestPawnAdvanceTwo(t *testing.T) {
	b := Board{}
	p := Pawn{}
	p.Team = TEAM_WHITE
	to := p.AvailableMoves(Location{FILE_B, RANK_2}, b)
	if !inLocations(Location{FILE_B, RANK_4}, to) {
		t.Fatal("Did not find move to advance two")
	}
	b.Squares[FILE_B][RANK_4] = &Pawn{}
	to = p.AvailableMoves(Location{FILE_B, RANK_2}, b)
	if inLocations(Location{FILE_B, RANK_4}, to) {
		t.Fatal("Found move to advance two when there was a piece there")
	}
	b.Squares[FILE_B][RANK_4] = nil
	b.Squares[FILE_B][RANK_3] = &Pawn{}
	to = p.AvailableMoves(Location{FILE_B, RANK_2}, b)
	if inLocations(Location{FILE_B, RANK_4}, to) {
		t.Fatal("Found move to advance two when there was a piece blocking")
	}
}

func TestPawnTakePiece(t *testing.T) {
	b := Board{}
	p := Pawn{}
	p.Team = TEAM_WHITE
	to := p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if inLocations(Location{FILE_A, RANK_5}, to) {
		t.Fatal("Found an attack to the left when it shouldn't have")
	}
	to = p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if inLocations(Location{FILE_C, RANK_5}, to) {
		t.Fatal("Found an attack to the right when it shouldn't have")
	}
	enemyPawn := &Pawn{}
	enemyPawn.Team = TEAM_BLACK
	b.Squares[FILE_A][RANK_5] = enemyPawn
	b.Squares[FILE_C][RANK_5] = enemyPawn
	to = p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if !inLocations(Location{FILE_A, RANK_5}, to) {
		t.Fatal("Could not find attack to the left")
	}
	if !inLocations(Location{FILE_C, RANK_5}, to) {
		t.Fatal("Could not find attack to the right")
	}
}

func TestEnPassant(t *testing.T) {
	b := Board{}
	p := Pawn{}
	p.Team = TEAM_WHITE
	to := p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if inLocations(Location{FILE_A, RANK_5}, to) {
		t.Fatal("Found an attack to the left when it shouldn't have")
	}
	to = p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if inLocations(Location{FILE_C, RANK_5}, to) {
		t.Fatal("Found an attack to the right when it shouldn't have")
	}
	enemyPawn := &Pawn{}
	enemyPawn.Team = TEAM_BLACK
	enemyPawn.MoveWasAdvanceTwo = true
	b.Squares[FILE_A][RANK_4] = enemyPawn
	b.Squares[FILE_C][RANK_4] = enemyPawn
	to = p.AvailableMoves(Location{FILE_B, RANK_4}, b)
	if !inLocations(Location{FILE_A, RANK_5}, to) {
		t.Fatal("Could not find attack to the left")
	}
	if !inLocations(Location{FILE_C, RANK_5}, to) {
		t.Fatal("Could not find attack to the right")
	}
}
