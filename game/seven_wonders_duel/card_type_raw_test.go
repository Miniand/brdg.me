package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardTypeRaw(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.Layout[lLen-1][0] = CardLumberYard
	g.Layout[lLen-1][1] = CardStable
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play lumber yard"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, fmt.Sprintf(
		"discard %s",
		Cards[g.Layout[lLen-1][2]].Name,
	)))
	mickCoins := g.PlayerCoins[helper.Mick]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play stable"))
	assert.Equal(t, mickCoins, g.PlayerCoins[helper.Mick], "stable wasn't free")
}
