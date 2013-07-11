package game

import (
	"testing"
)

func TestRawCollection(t *testing.T) {
	for id, g := range RawCollection() {
		if id != g.Identifier() {
			t.Error("Got wrong id for game, expecting", id, "got",
				g.Identifier())
		}
	}
}

func TestCollection(t *testing.T) {
	for id, f := range Collection() {
		g, _ := f([]string{})
		if id != g.Identifier() {
			t.Error("Got wrong id for game, expecting", id, "got",
				g.Identifier())
		}
	}
}
