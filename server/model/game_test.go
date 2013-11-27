package model

import (
	"github.com/Miniand/brdg.me/game"
	"testing"
)

func TestGameSavingAndLoading(t *testing.T) {
	if modelTestShouldRun() {
		cleanTestingDatabase()
		g, err := game.Collection()["acquire"]([]string{"Mick", "Steve", "BJ"})
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
		t.Logf("%#v", loadedGm)
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
		if len(pl) != 3 || pl[0] != "Mick" || pl[1] != "Steve" ||
			pl[2] != "BJ" {
			t.Error("Players in loaded game don't match original")
			return
		}
	}
}
