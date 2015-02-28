package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestFreeBuild(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Hands[Mick][0] = Cards[CardPalace]

	// Palace is so expensive, also don't have the right wonder stage.
	assert.Error(t, cmd(g, Mick, "build 1"))
	assert.Error(t, cmd(g, Mick, "free 1"))

	g.Cards[Mick] = card.Deck{
		Cards[WonderStageOlympiaA2],
	}

	// Should be able to free build now
	assert.NoError(t, cmd(g, Mick, "free 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))

	// Shouldn't be able to free build for the rest of the round
	for g.Round == 1 {
		assert.Error(t, cmd(g, Mick, "free 1"))
		for p := range g.Players {
			assert.NoError(t, cmd(g, p, "discard 1"))
		}
	}

	// New round, should be able to free build again
	assert.NoError(t, cmd(g, Mick, "free 1"))
}
