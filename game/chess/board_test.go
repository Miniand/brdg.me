package chess

import (
	"fmt"
	"github.com/beefsack/brdg.me/render"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	b := Board{}
	if !b.IsEmpty(Location{FILE_A, RANK_1}) {
		t.Fatal("A1 should be empty")
	}
	b.Squares[FILE_A][RANK_1] = &Pawn{}
	if b.IsEmpty(Location{FILE_A, RANK_1}) {
		t.Fatal("A1 should not be empty")
	}
}

func TestRender(t *testing.T) {
	b := Board{}
	b.Squares[FILE_C][RANK_5] = &Pawn{Piece{TEAM_BLACK}}
	output, _ := render.RenderTerminal(b.Render())
	fmt.Println(output)
}
