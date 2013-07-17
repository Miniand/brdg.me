package main

import (
	"github.com/beefsack/brdg.me/game"
	"os"
	"testing"
)

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DATABASE") != ""
}

func TestSavingAndLoading(t *testing.T) {
	if modelTestShouldRun() {
		g, err := game.Collection()["tic_tac_toe"]([]string{"Mick", "Steve"})
		if err != nil {
			t.Error(err)
			return
		}
		gm, err := SaveGame(g)
		if err != nil {
			t.Error(err)
			return
		}
		loadedGm, err := LoadGame(gm.Id)
		if err != nil {
			t.Error(err)
			return
		}
		loadedG, err := loadedGm.ToGame()
		if err != nil {
			t.Error(err)
			return
		}
		if gm.Id != loadedGm.Id {
			t.Error("Id doesn't match")
			return
		}
		pl := loadedG.PlayerList()
		if len(pl) != 2 || pl[0] != "Mick" || pl[1] != "Steve" {
			t.Error("Players in loaded game don't match original")
			return
		}
	}
}
