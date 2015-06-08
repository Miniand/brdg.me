package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestTakeCommand(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Hands[Mick][0] = Cards[WonderStageHalicarnassusA2]
	g.Cards[Mick] = card.Deck{
		Cards[CardOreVein],
		Cards[CardFoundry],
	}
	g.Discard = card.Deck{
		Cards[CardPalace],
	}

	assert.NoError(t, cmd(g, Mick, "build 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))

	assert.Equal(t, []string{players[Mick]}, g.WhoseTurn())

	assert.NoError(t, cmd(g, Mick, "take 1"))

	assert.Len(t, g.Cards[Mick], 4)
	assert.Equal(t, CardPalace, g.Cards[Mick][3].(Carder).GetCard().Name)
	assert.Len(t, g.Discard, 2)
	assert.Equal(t, players, g.WhoseTurn())
}

func TestTakeCommand_currentlyDiscarded(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Hands[Mick][0] = Cards[WonderStageHalicarnassusA2]
	g.Cards[Mick] = card.Deck{
		Cards[CardOreVein],
		Cards[CardFoundry],
	}

	assert.NoError(t, cmd(g, Mick, "build 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))

	assert.Equal(t, []string{players[Mick]}, g.WhoseTurn())
}

func TestTakeCommand_empty(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Hands[Mick][0] = Cards[WonderStageHalicarnassusA2]
	g.Hands[Steve][0] = Cards[CardLumberYard]
	g.Hands[Greg][0] = Cards[CardStonePit]
	g.Cards[Mick] = card.Deck{
		Cards[CardOreVein],
		Cards[CardFoundry],
	}
	g.Discard = card.Deck{
		Cards[CardOreVein], // Not buildable so should be ignored.
	}

	assert.NoError(t, cmd(g, Mick, "build 1"))
	assert.NoError(t, cmd(g, Steve, "build 1"))
	assert.NoError(t, cmd(g, Greg, "build 1"))

	assert.Equal(t, players, g.WhoseTurn())
}

func TestTakeCommand_alreadyBuild(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Hands[Mick][0] = Cards[WonderStageHalicarnassusA2]
	g.Cards[Mick] = card.Deck{
		Cards[CardOreVein],
		Cards[CardFoundry],
	}
	g.Discard = card.Deck{
		Cards[CardPalace],
		Cards[CardOreVein],
	}

	assert.NoError(t, cmd(g, Mick, "build 1"))
	assert.NoError(t, cmd(g, Steve, "discard 1"))
	assert.NoError(t, cmd(g, Greg, "discard 1"))

	assert.Equal(t, []string{players[Mick]}, g.WhoseTurn())

	assert.Error(t, cmd(g, Mick, "take 2"))
	assert.NoError(t, cmd(g, Mick, "take 1"))
}
