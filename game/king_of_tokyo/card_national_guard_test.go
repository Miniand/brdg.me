package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardNationalGuard(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.FaceUpCards = []int{NationalGuard}
	g.CurrentRoll = []int{DieEnergy, DieEnergy, DieEnergy}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	initialHealth := g.Boards[helper.Mick].Health
	initialVP := g.Boards[helper.Mick].VP
	assert.NoError(t, helper.Cmd(g, helper.Mick, "buy national guard"))
	assert.Equal(t, initialHealth-2, g.Boards[helper.Mick].Health)
	assert.Equal(t, initialVP+2, g.Boards[helper.Mick].VP)
}
