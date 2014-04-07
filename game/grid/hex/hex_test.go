package hex

import (
	"reflect"
	"testing"

	"github.com/Miniand/brdg.me/game/grid"
)

func TestSetTile(t *testing.T) {
	g := Grid{}
	g.SetTile(grid.Loc{-5, 4}, "egg")
	if g.Tile(grid.Loc{-5, 4}) != "egg" {
		t.Fatal("Did not get egg")
	}
}

func TestNeighbours(t *testing.T) {
	hex := Grid{}
	given := []grid.Loc{
		grid.Loc{2, 0},
	}
	expected := [][]grid.Loc{
		[]grid.Loc{
			grid.Loc{2, -1},
			grid.Loc{3, -1},
			grid.Loc{3, 0},
			grid.Loc{2, 1},
			grid.Loc{1, 0},
			grid.Loc{1, -1},
		},
	}
	for i, g := range given {
		actual := hex.Neighbours(g)
		if len(actual) != len(expected[i]) {
			t.Fatal("Lengths do not match")
		}
		for i2, a := range actual {
			if a.X != expected[i][i2].X || a.Y != expected[i][i2].Y {
				t.Fatalf("%#v did not match expected %#v", a, expected[i][i2])
			}
		}
	}
}

func TestBounds(t *testing.T) {
	hex := Grid{}
	hex.SetTile(grid.Loc{-6, 0}, "blah")
	hex.SetTile(grid.Loc{0, -7}, "blah")
	hex.SetTile(grid.Loc{8, 0}, "blah")
	hex.SetTile(grid.Loc{0, 9}, "blah")
	lower, upper := hex.Bounds()
	if lower.X != -6 || lower.Y != -7 || upper.X != 8 || upper.Y != 9 {
		t.Fatalf("Bounds did not match, got: %#v %#v", lower, upper)
	}
}

func TestFind(t *testing.T) {
	hex := Grid{}
	hex.SetTile(grid.Loc{-4, 6}, "moo")
	hex.SetTile(grid.Loc{-3, 6}, "bag")
	hex.SetTile(grid.Loc{-2, 6}, "mazomba")
	found, ok := hex.Find("bag")
	if !ok {
		t.Fatal("Could not find bag")
	}
	if !reflect.DeepEqual(found, grid.Loc{-3, 6}) {
		t.Fatal("Could not find at -3, 6")
	}
	_, ok = hex.Find("should not exist")
	if ok {
		t.Fatal("Found something that should not exist")
	}
}
