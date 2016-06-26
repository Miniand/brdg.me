package seven_wonders_duel

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestPlayCommandFreeBuild(t *testing.T) {
	g := atFirstRound(t)
	lLen := len(g.Layout)
	g.Layout[lLen-1][0] = CardStatue

	g.PlayerCards[helper.Mick] = []int{CardTheater}

	mickCoin := g.PlayerCoins[helper.Mick]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play statue"))
	assert.Equal(t, mickCoin, g.PlayerCoins[helper.Mick])
}
