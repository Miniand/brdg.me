package texas_holdem

import (
	"testing"
)

func mockGame(t *testing.T) *Game {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve"})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestStart(t *testing.T) {
	g := mockGame(t)
	if len(g.Players) == 0 {
		t.Fatal("Didn't set players")
	}
}
