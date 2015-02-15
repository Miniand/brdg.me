package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestCardPlayFinalCard_With(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	// Trim down to two for testing purposes.
	for h := range g.Hands {
		g.Hands[h] = g.Hands[h][:2]
	}
	g.Cards[Steve] = card.Deck{
		Cards[WonderStageBabylonB2],
	}
	assert.NoError(t, cmd(g, Mick, "discard 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))
	assert.Equal(t, g.Round, 1)
	assert.Equal(t, g.WhoseTurn(), []string{players[Steve]})
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.Equal(t, g.Round, 2)
	assert.Equal(t, g.WhoseTurn(), players)
}

func TestCardPlayFinalCard_Without(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	// Trim down to two for testing purposes.
	for h := range g.Hands {
		g.Hands[h] = g.Hands[h][:2]
	}
	assert.NoError(t, cmd(g, Mick, "discard 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))
	assert.Equal(t, g.Round, 2)
	assert.Equal(t, g.WhoseTurn(), players)
}
