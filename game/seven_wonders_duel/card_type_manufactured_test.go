package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardTypeManufactured(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.Layout[lLen-1][0] = CardGlassworks
	g.Layout[lLen-1][1] = CardApothecary
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play glassworks"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, fmt.Sprintf(
		"discard %s",
		Cards[g.Layout[lLen-1][2]].Name,
	)))
	mickCoins := g.PlayerCoins[helper.Mick]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play apothecary"))
	assert.Equal(t, mickCoins, g.PlayerCoins[helper.Mick], "apothecary wasn't free")
}
