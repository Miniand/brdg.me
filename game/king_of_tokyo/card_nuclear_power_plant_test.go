package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardNuclearPowerPlant(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.FaceUpCards = []int{NuclearPowerPlant}
	g.Boards[helper.Mick].Health = 3
	g.Boards[helper.Mick].Energy = 6
	g.CurrentRoll = []int{DieEnergy}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "buy nuclear"))
	assert.Equal(t, 6, g.Boards[helper.Mick].Health)
	assert.Equal(t, 2, g.Boards[helper.Mick].VP)
}
