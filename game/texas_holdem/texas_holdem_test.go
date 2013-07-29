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

func TestNextPhaseOnInitialFold(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "BJ"})
	if err != nil {
		t.Fatal(err)
	}
	// First player folds
	err = g.PlayerAction(g.Players[g.CurrentPlayer], "fold", []string{})
	if err != nil {
		t.Fatal(err)
	}
	// Next two players call and check, should flop
	err = g.PlayerAction(g.Players[g.CurrentPlayer], "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	err = g.PlayerAction(g.Players[g.CurrentPlayer], "check", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("No flop, community cards:", g.CommunityCards)
	}
}
