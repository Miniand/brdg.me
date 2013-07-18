package lost_cities

import (
	"testing"
)

func TestStart(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Error(err)
	}
	if game.PlayerList()[0] != "Mick" {
		t.Error("Mick is not a player")
	}
	if game.PlayerList()[1] != "Steve" {
		t.Error("Steve is not a player")
	}
}
