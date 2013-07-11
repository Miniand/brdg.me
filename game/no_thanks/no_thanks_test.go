package no_thanks

import (
	"sort"
	"testing"
)

func TestStart(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	if g.Players[0] != "Mick" || g.Players[1] != "Steve" {
		t.Error("Players aren't Mick and Steve, got", g.Players)
		return
	}
	if g.CurrentlyMoving != "Mick" && g.CurrentlyMoving != "Steve" &&
		g.CurrentlyMoving != "Barabbas" {
		t.Error("Currently moving not set to Mick or Steve, got",
			g.CurrentlyMoving)
		return
	}
}

func TestAllCards(t *testing.T) {
	g := &Game{}
	cards := g.AllCards()
	if len(cards) != 33 {
		t.Error("There weren't 33 cards, got", len(cards))
		return
	}
	if cards[0] != 3 {
		t.Error("Expected the first card to be 3, got", cards[0])
		return
	}
	if cards[32] != 35 {
		t.Error("Expected the thirty third card to be 35, got", cards[32])
		return
	}
}

func TestInitCards(t *testing.T) {
	g := &Game{}
	g.InitCards()
	if len(g.RemainingCards) != 24 {
		t.Error("Expected there to be 24 cards in the stack, got",
			len(g.RemainingCards))
		return
	}
	for _, c := range g.RemainingCards {
		if c < 3 || c > 35 {
			t.Error("Expected cards to be between 3 and 35, got", c)
			return
		}
	}
}

func TestInitPlayerChips(t *testing.T) {
	g := &Game{}
	g.InitPlayerChips()
	for _, p := range g.Players {
		if g.PlayerChips[p] != 11 {
			t.Error("Expected player chips to be 11, got", g.PlayerChips[p])
			return
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	data, err := g.Encode()
	if err != nil {
		t.Error(err)
		return
	}
	g2 := &Game{}
	err = g2.Decode(data)
	if err != nil {
		t.Error(err)
		return
	}
	if g.Players[0] != g2.Players[0] && g.Players[1] != g2.Players[1] {
		t.Error("Players don't match", g.Players, g2.Players)
		return
	}
}

func TestAssertTurn(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.CurrentlyMoving = "Steve"
	if g.AssertTurn("Mick") == nil {
		t.Error("Checking if it's Mick's turn should error")
		return
	}
	if g.AssertTurn("Steve") != nil {
		t.Error("Checking if it's Steve's turn should not error")
		return
	}
	g.RemainingCards = []int{}
	if g.AssertTurn("Steve") == nil {
		t.Error("Checking if it's Steve's turn should error when finished")
		return
	}
}

func TestIsFinished(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	if g.IsFinished() {
		t.Error("Game should not be finished immediately after starting it")
		return
	}
	g.RemainingCards = []int{}
	if !g.IsFinished() {
		t.Error("Game should be finished when there are no cards left")
		return
	}
}

func TestPass(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	initialPlayer := g.CurrentlyMoving
	initialCardCount := len(g.RemainingCards)
	initialCentreChips := g.CentreChips
	initialPlayerChips := g.PlayerChips[initialPlayer]
	err = g.Pass(initialPlayer)
	if err != nil {
		t.Error(err)
		return
	}
	if len(g.RemainingCards) != initialCardCount {
		t.Error("The card count changed when it shouldn't have")
		return
	}
	if g.CentreChips != initialCentreChips+1 {
		t.Error("Centre chips didn't increase by 1")
		return
	}
	if g.PlayerChips[initialPlayer] != initialPlayerChips-1 {
		t.Error("Expected player chips to be reduced by 1 but it wasn't")
		return
	}
	if g.CurrentlyMoving == initialPlayer {
		t.Error("Didn't move to the next player")
		return
	}
}

func TestTake(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	initialPlayer := g.CurrentlyMoving
	initialCardCount := len(g.RemainingCards)
	g.CentreChips = 5 // Set the centre chips for the sake of testing
	initialCentreChips := g.CentreChips
	initialPlayerChips := g.PlayerChips[initialPlayer]
	topCard := g.PeekTopCard()
	err = g.Take(initialPlayer)
	if err != nil {
		t.Error(err)
		return
	}
	if len(g.RemainingCards) != initialCardCount-1 {
		t.Error("The card count didn't reduce by 1")
		return
	}
	if g.CentreChips != 0 {
		t.Error("Centre chips should have been 0 after taking")
		return
	}
	if g.PlayerChips[initialPlayer] != initialPlayerChips+initialCentreChips {
		t.Error("Player didn't take the centre chips")
		return
	}
	if len(g.PlayerHands[initialPlayer]) != 1 ||
		g.PlayerHands[initialPlayer][0] != topCard {
		t.Error("Player didn't take the top card into their hand")
		return
	}
	if g.CurrentlyMoving != initialPlayer {
		t.Error("Moved to next player when it shouldn't have")
		return
	}
}

func TestPlayerHandSorted(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Mick"] = []int{5, 3, 6, 4, 87}
	if !sort.IntsAreSorted(g.PlayerHandSorted("Mick")) {
		t.Error("Mick's hand wasn't sorted when fetching via sorted method")
		return
	}
}

func TestPlayerHandGrouped(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Mick"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	grouping := g.PlayerHandGrouped("Mick")
	if len(grouping) != 4 {
		t.Error("Expected 4 groups, got", len(grouping))
		return
	}
	if len(grouping[0]) != 1 {
		t.Error("Expected group 1 to be 1 card, got", len(grouping[0]))
		return
	}
	if grouping[0][0] != 3 {
		t.Error("Expected group 1 to be [3], got", grouping[0])
		return
	}
	if len(grouping[1]) != 2 {
		t.Error("Expected group 2 to be 2 cards, got", len(grouping[1]))
		return
	}
	if grouping[1][0] != 5 || grouping[1][1] != 6 {
		t.Error("Expected group 2 to be [5 6], got", grouping[1])
		return
	}
	if len(grouping[2]) != 3 {
		t.Error("Expected group 3 to be 3 cards, got", len(grouping[2]))
		return
	}
	if grouping[2][0] != 8 || grouping[2][1] != 9 || grouping[2][2] != 10 {
		t.Error("Expected group 3 to be [8 9 10], got", grouping[2])
		return
	}
	if len(grouping[3]) != 2 {
		t.Error("Expected group 4 to be 2 cards, got", len(grouping[3]))
		return
	}
	if grouping[3][0] != 15 || grouping[3][1] != 16 {
		t.Error("Expected group 4 to be [15 16], got", grouping[3])
		return
	}
}

func TestPlayerHandScore(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Mick"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	g.PlayerChips["Mick"] = 10
	expectedScore := 3 + 5 + 8 + 15
	if g.PlayerHandScore("Mick") != expectedScore {
		t.Error("Expected score of", expectedScore, "got",
			g.PlayerHandScore("Mick"))
		return
	}
}

func TestFinalPlayerScore(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Mick"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	g.PlayerChips["Mick"] = 10
	expectedScore := 3 + 5 + 8 + 15 - 10
	if g.FinalPlayerScore("Mick") != expectedScore {
		t.Error("Expected score of", expectedScore, "got",
			g.FinalPlayerScore("Mick"))
		return
	}
}

func TestWhoseTurn(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	if len(g.WhoseTurn()) != 1 || g.WhoseTurn()[0] != g.CurrentlyMoving {
		t.Error("Expected turn to be", g.CurrentlyMoving, "got", g.WhoseTurn())
		return
	}
}

func TestWinners(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Barabbas", "Mick", "Steve"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Barabbas"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	g.PlayerChips["Barabbas"] = 3
	g.PlayerHands["Mick"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	g.PlayerChips["Mick"] = 10
	g.PlayerHands["Steve"] = []int{5, 8, 3, 10, 9, 6, 16, 17}
	g.PlayerChips["Steve"] = 11
	g.RemainingCards = []int{}
	winners := g.Winners()
	if len(winners) != 2 || winners[0] != "Mick" || winners[1] != "Steve" {
		t.Error("Expected winners to be Mick and Steve, got", winners)
		return
	}
}

func TestPlayerActions(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.CurrentlyMoving = "Steve"
	topCard := g.PeekTopCard()
	err = g.PlayerAction("Steve", "Pass", []string{})
	if err != nil {
		t.Error(err)
		return
	}
	if g.PlayerChips["Steve"] != 10 {
		t.Error("Expected Steve's chips to be 10, got", g.PlayerChips["Steve"])
		return
	}
	err = g.PlayerAction("Barabbas", "taKE", []string{})
	if err != nil {
		t.Error(err)
		return
	}
	if len(g.PlayerHands["Barabbas"]) != 1 ||
		g.PlayerHands["Barabbas"][0] != topCard {
		t.Error("Barabbas did not take the card, hand:", g.PlayerHands)
		return
	}
	if g.PlayerChips["Barabbas"] != 12 {
		t.Error("Expected Barabbas's chips to be 12, got",
			g.PlayerChips["Barabbas"])
		return
	}
}

func TestRenderCardsForPlayer(t *testing.T) {
	g := &Game{}
	err := g.Start([]string{"Mick", "Steve", "Barabbas"})
	if err != nil {
		t.Error(err)
		return
	}
	g.PlayerHands["Mick"] = []int{5, 8, 3, 10, 9, 15, 6, 16}
	g.PlayerChips["Mick"] = 10
	expectedOutput := "3   5 {{b}}6{{_b}}   {{b}}8{{_b}} 9 10   15 16"
	output := g.RenderCardsForPlayer("Mick", 7)
	if output != expectedOutput {
		t.Error("Expected output", expectedOutput, "but got", output)
		return
	}
}
