package jaipur

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestSellCommand_Call(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	g.Hands[helper.Mick] = []int{GoodGold}
	assert.Error(t, helper.Cmd(g, helper.Mick, "sell 2 gold"))
	assert.Error(t, helper.Cmd(g, helper.Mick, "sell 1 gold"))
	g.Hands[helper.Mick] = []int{GoodGold, GoodLeather, GoodGold}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "sell 2 gold"))
	assert.Equal(t, []int{6, 6}, g.Tokens[helper.Mick])
	assert.Equal(t, []int{5, 5, 5}, g.Goods[GoodGold])
	assert.Equal(t, []int{GoodLeather}, g.Hands[helper.Mick])
}
