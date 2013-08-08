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
	if len(g.CommunityCards) != 0 {
		t.Fatal("Cards were already drawn:", g.CommunityCards)
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
	g.FirstBettingPlayer = 2
	g.Bets = []int{
		0: 5,
		1: 10,
		2: 0,
	}
	err = g.PlayerAction("Mick", "call", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 0 {
		t.Fatal("Flopped too early")
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

func TestAllInAboveOtherPlayer(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"BJ", "Mick"})
	if err != nil {
		t.Fatal(err)
	}
	g.CommunityCards, g.Deck = g.Deck.PopN(3)
	g.CurrentDealer = 0
	g.CurrentPlayer = 0
	g.FirstBettingPlayer = 0
	g.Bets = []int{
		0: 5,
		1: 10,
	}
	g.PlayerMoney = []int{
		0: 10,
		1: 20,
	}
	// Go all in with BJ
	err = g.PlayerAction("BJ", "allin", []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.CommunityCards) != 3 || g.CurrentPlayer != 1 {
		t.Fatal("Game progressed without letting Mick call")
	}
}

func TestAllPlayersAllInWhenBlindsBiggerThanCash(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"BJ", "Mick"})
	if err != nil {
		t.Fatal(err)
	}
	g.PlayerMoney = []int{
		0: 3,
		1: 3,
	}
	g.NewHand()
	if !g.IsFinished() {
		t.Fatal("Game didn't finish when players had lower money than blinds")
	}
}

func TestNextPlayerIsSkippedOnNextPhaseWhenNoMoney(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"BJ", "Mick", "Steve"})
	if err != nil {
		t.Fatal(err)
	}
	g.CurrentDealer = 0
	g.CurrentPlayer = 0
	g.FirstBettingPlayer = 0
	g.Bets = []int{
		0: 10,
		1: 6,
		2: 10,
	}
	g.PlayerMoney = []int{
		0: 10,
		1: 0,
		2: 100,
	}
	// Skip to next phase manually
	g.NextPhase()
	if g.CurrentPlayer != 2 {
		t.Fatal("Didn't skip over Mick on new phase even though he is all in, got:", g.CurrentPlayer)
	}
}
