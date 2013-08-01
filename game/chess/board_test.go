package chess

import (
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
	b := InitialBoard()
	expected := `  a b c d e f g h
8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 8
7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 7
6 · · · · · · · · 6
5 · · · · · · · · 5
4 · · · · · · · · 4
3 · · · · · · · · 3
2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ 2
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 1
  a b c d e f g h`
	output, _ := render.RenderPlain(b.Render())
	if output != expected {
		t.Fatalf("Board was not:\n%s\nGot\n%s", expected, output)
	}
}
