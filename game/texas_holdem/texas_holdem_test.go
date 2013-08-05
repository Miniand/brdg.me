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
	if len(g.CommunityCards) != 0 {
		t.Fatal("Cards were already drawn:", g.CommunityCards)
	}
	err = g.PlayerAction(g.Players[g.CurrentPlayer], "check", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("No flop, community cards:", g.CommunityCards)
	}
}

func TestDealerRaiseWhenLastPlayer(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"BJ", "Pete", "Mick"})
	if err != nil {
		t.Fatal(err)
	}
	g.CurrentDealer = 2
	g.CurrentPlayer = 2
	g.LastRaisingPlayer = 2
	g.Bets = []int{
		0: 5,
		1: 10,
		2: 0,
	}
	err = g.PlayerAction("Mick", "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	err = g.PlayerAction("BJ", "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 0 {
		t.Fatal("Flopped too early")
	}
	err = g.PlayerAction("Pete", "check", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("Flop didn't happen")
	}
	err = g.PlayerAction("BJ", "check", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("Turn happened too early")
	}
	err = g.PlayerAction("Pete", "check", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("Turn happened too early")
	}
	err = g.PlayerAction("Mick", "raise", []string{"10"})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("Turn happened too early")
	}
	err = g.PlayerAction("BJ", "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 {
		t.Fatal("Turn happened too early")
	}
	err = g.PlayerAction("Pete", "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 4 {
		t.Fatal("Turn didn't happen")
	}
}
