package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardParasiticTentacles(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.FaceUpCards = []int{ParasiticTentacles}
	g.Boards[helper.Steve].Cards = []int{Jets, NationalGuard}
	g.Boards[helper.Mick].Energy = 6
	g.CurrentRoll = []int{DieEnergy, DieEnergy, DieEnergy}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	assert.Error(t, helper.Cmd(g, helper.Mick, "buy jets"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "buy paris"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "buy jets"))
	assert.Equal(t, 0, g.Boards[helper.Mick].Energy)
	assert.Equal(t, 5, g.Boards[helper.Steve].Energy)
	assert.Contains(t, g.Boards[helper.Mick].Cards, Jets)
	assert.NotContains(t, g.Boards[helper.Steve].Cards, Jets)
}
