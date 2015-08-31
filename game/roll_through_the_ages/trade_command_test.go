package roll_through_the_ages

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestTradeCommand(t *testing.T) {
	g := &Game{}
	g.Start(helper.Players[:3])
	g.CurrentPlayer = helper.Mick
	g.Phase = PhaseRoll
	g.RolledDice = []int{
		DiceFood,
		DiceFood,
		DiceFood,
	}
	g.KeptDice = []int{}
	g.Boards[helper.Mick].Developments[DevelopmentEngineering] = true
	g.Boards[helper.Mick].Goods[GoodStone] = 3
	assert.NoError(t, helper.Cmd(g, helper.Mick, "next"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "trade 3"))
	assert.Equal(t, 0, g.Boards[helper.Mick].Goods[GoodStone])
	assert.NoError(t, helper.Cmd(g, helper.Mick, "build 9 great"))
	assert.Equal(t, 9, g.Boards[helper.Mick].Monuments[MonumentGreatPyramid])
}
