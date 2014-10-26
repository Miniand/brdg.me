package chess

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Miniand/brdg.me/render"
)

func parseBoard(input string, t *testing.T) (b Board) {
	ranks := strings.Split(strings.TrimSpace(input), "\n")
	if len(ranks) != 8 {
		t.Fatal("Not 8 ranks in input")
	}
	for rank, squares := range ranks {
		runes := bytes.Runes([]byte(squares))
		if len(runes) != 8 {
			t.Fatalf("Not 8 runes in row %d", rank)
		}
		for file, contents := range runes {
			l := Location{file, rank}
			if IsValidLocation(l) {
				var piece Piecer
				switch contents {
				case '♜':
					p := &Rook{}
					p.Team = TEAM_BLACK
					piece = p
				case '♞':
					p := &Knight{}
					p.Team = TEAM_BLACK
					piece = p
				case '♝':
					p := &Bishop{}
					p.Team = TEAM_BLACK
					piece = p
				case '♛':
					p := &Queen{}
					p.Team = TEAM_BLACK
					piece = p
				case '♚':
					p := &King{}
					p.Team = TEAM_BLACK
					piece = p
				case '♟':
					p := &Pawn{}
					p.Team = TEAM_BLACK
					piece = p
				case '♖':
					p := &Rook{}
					p.Team = TEAM_WHITE
					piece = p
				case '♘':
					p := &Knight{}
					p.Team = TEAM_WHITE
					piece = p
				case '♗':
					p := &Bishop{}
					p.Team = TEAM_WHITE
					piece = p
				case '♕':
					p := &Queen{}
					p.Team = TEAM_WHITE
					piece = p
				case '♔':
					p := &King{}
					p.Team = TEAM_WHITE
					piece = p
				case '♙':
					p := &Pawn{}
					p.Team = TEAM_WHITE
					piece = p
				}
				// Note that we need the inverse of the rank
				b.Squares[file][RANK_8-rank] = piece
			}
		}
	}
	return
}

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
	output := render.RenderPlain(b.Render())
	if output != expected {
		t.Fatalf("Board was not:\n%s\nGot\n%s", expected, output)
	}
}
