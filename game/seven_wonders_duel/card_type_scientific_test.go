package seven_wonders_duel

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardTypeScientificProgress(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.PlayerCoins[helper.Mick] = 10
	g.Layout[lLen-1][0] = CardScriptorium
	g.Layout[lLen-1][1] = CardPharmacist
	g.Layout[lLen-1][2] = CardLibrary
	// Some cards for Steve to discard
	g.Layout[lLen-1][3] = CardTheater
	g.Layout[lLen-1][4] = CardAltar

	g.ProgressTokens[0] = ProgressAgriculture

	assert.NoError(t, helper.Cmd(g, helper.Mick, "play scriptorium"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "discard theater"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play pharmacist"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "discard altar"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play library"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "choose agriculture"))
	assert.Contains(t, g.PlayerCards[helper.Mick], ProgressAgriculture)
	assert.Equal(t, helper.Steve, g.CurrentPlayer)
}

func TestCardTypeScientificProgressNoneLeft(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.PlayerCoins[helper.Mick] = 10
	g.Layout[lLen-1][0] = CardScriptorium
	g.Layout[lLen-1][1] = CardPharmacist
	g.Layout[lLen-1][2] = CardLibrary
	// Some cards for Steve to discard
	g.Layout[lLen-1][3] = CardTheater
	g.Layout[lLen-1][4] = CardAltar

	g.ProgressTokens = []int{}

	assert.NoError(t, helper.Cmd(g, helper.Mick, "play scriptorium"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "discard theater"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play pharmacist"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "discard altar"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play library"))
	assert.Equal(t, helper.Steve, g.CurrentPlayer)
}

func TestCardTypeScientificVictory(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.PlayerCards[helper.Mick] = []int{
		CardLibrary,
		CardAcademy,
		CardLaboratory,
		ProgressLaw,
	}
	g.Layout[lLen-1][0] = CardObservatory
	g.Layout[lLen-1][1] = CardPharmacist
	// Some cards for Steve to discard
	g.Layout[lLen-1][2] = CardTheater

	assert.NoError(t, helper.Cmd(g, helper.Mick, "play observatory"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "discard theater"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play pharmacist"))
	assert.True(t, g.IsFinished())
	assert.Equal(t, []string{g.Players[helper.Mick]}, g.Winners())
}
