package liars_dice

import (
	"github.com/Miniand/brdg.me/command"
	"testing"
)

func TestStart(t *testing.T) {
	g := &Game{}
	p := []string{"Mick", "Steve", "BJ"}
	if err := g.Start(p); err != nil {
		t.Fatal(err)
	}
	pl := g.PlayerList()
	if len(pl) != len(p) || pl[0] != p[0] || pl[1] != p[1] || pl[2] != p[2] {
		t.Fatal("PlayerList doesn't match players that started the game")
	}
	wt := g.WhoseTurn()
	if len(wt) != 1 || (wt[0] != p[0] && wt[0] != p[1] && wt[0] != p[2]) {
		t.Fatal("WhoseTurn doesn't report one of the specified players")
	}
	if len(g.PlayerDice) != len(p) {
		t.Fatal("PlayerDice hasn't been initialised for each player")
	}
	for i, _ := range g.Players {
		if len(g.PlayerDice[i]) != START_DICE_COUNT {
			t.Fatalf("PlayerDice for %s has not been initialised to 5 dice",
				g.Players[i])
		}
		for _, d := range g.PlayerDice[i] {
			if d < 1 || d > 6 {
				t.Fatalf("Dice isn't in the range of 1 to 6")
			}
		}
	}
}

func TestExampleRound(t *testing.T) {
	g := &Game{}
	p := []string{"Mick", "Steve", "BJ"}
	if err := g.Start(p); err != nil {
		t.Fatal(err)
	}
	// Override the game values so we know what's there
	g.PlayerDice = [][]int{
		// Mick
		[]int{1, 3, 4, 4, 6},
		// Steve
		[]int{2, 2, 3, 3, 3},
		// BJ
		[]int{1},
	}
	// First player is mick
	g.CurrentPlayer = 0
	// Make sure we can't call on the first turn
	if _, err := command.CallInCommands("Mick", g, "call",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when calling on the first turn")
	}
	// Start with a few legit commands
	if _, err := command.CallInCommands("Mick", g, "bid 2 5",
		g.Commands()); err != nil {
		t.Fatal(err)
	}
	if _, err := command.CallInCommands("Steve", g, "bid 2 6",
		g.Commands()); err != nil {
		t.Fatal(err)
	}
	if _, err := command.CallInCommands("BJ", g, "bid 3 5",
		g.Commands()); err != nil {
		t.Fatal(err)
	}
	// Do a few illegal commands and make sure they're picked up
	if _, err := command.CallInCommands("Mick", g, "bid 3 5",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when making same bid")
	}
	if _, err := command.CallInCommands("Mick", g, "bid 3 3",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when bidding a lower value dice")
	}
	if _, err := command.CallInCommands("Mick", g, "bid 2 6",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when reducing the quantity")
	}
	if _, err := command.CallInCommands("Mick", g, "bid 3 7",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when making an bid of an invalid dice value")
	}
	if _, err := command.CallInCommands("BJ", g, "bid 6 5",
		g.Commands()); err == nil {
		t.Fatal("Didn't fail when BJ barged in")
	}
	// Call it and check
	if _, err := command.CallInCommands("Mick", g, "call",
		g.Commands()); err != nil {
		t.Fatal(err)
	}
	if len(g.PlayerDice[2]) != 0 {
		t.Fatal("BJ should have lost his dice")
	}
	if len(g.PlayerDice[0]) != 5 && len(g.PlayerDice[1]) != 5 {
		t.Fatal("Mick and Steve shouldn't have lost dice")
	}
	if g.CurrentPlayer != 1 {
		t.Fatal("Steve didn't become the current player")
	}
	if len(g.ActivePlayers()) != 2 {
		t.Fatal("BJ wasn't eliminated")
	}
}

func TestPlayerElimination(t *testing.T) {
	g := &Game{}
	p := []string{"Mick", "Steve", "BJ", "Ross"}
	if err := g.Start(p); err != nil {
		t.Fatal(err)
	}
	g.PlayerDice[0] = []int{}
	g.PlayerDice[2] = []int{}
	eliminated := g.EliminatedPlayerList()
	if len(eliminated) != 2 {
		t.Fatal("Two players weren't eliminated, got:", eliminated)
	}
	if eliminated[0] != "Mick" {
		t.Fatal("Mick was not eliminated, got:", eliminated)
	}
	if eliminated[1] != "BJ" {
		t.Fatal("BJ was not eliminated, got:", eliminated)
	}
}
