package seven_wonders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCommercial(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Hands[Steve][0] = Cards[CardTavern]
	origMoney := g.Coins[Steve]
	assert.NoError(t, cmd(g, Mick, "discard 1"))
	assert.NoError(t, cmd(g, Steve, "build 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))
	assert.Equal(t, origMoney+5, g.Coins[Steve])
}
