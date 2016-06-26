package seven_wonders_duel

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardTypeCivilian(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.Layout[lLen-1][0] = CardTheater
	mickVP := g.PlayerVP(helper.Mick)
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play theater"))
	assert.Equal(t, mickVP+3, g.PlayerVP(helper.Mick), "theater didn't provide 3 VP")
}
