package liars_dice

import (
	"testing"
)

func TestStart(t *testing.T) {
	g := &Game{}
	p := []string{"Mick", "Steve", "BJ"}
	err := g.Start(p)
	if err != nil {
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
