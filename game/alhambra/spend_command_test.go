package alhambra

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestSpendCommand_multipleSameCard(t *testing.T) {
	g := &Game{}
	g.Start(helper.Players)
	g.CurrentPlayer = Mick
	g.Boards[Mick].Cards = card.Deck{
		Card{CurrencyBlue, 1},
		Card{CurrencyBlue, 1},
		Card{CurrencyBlue, 1},
	}
	g.Tiles[CurrencyBlue] = NewTile(TileTypeTower, 3)

	assert.Error(t, helper.Cmd(g, helper.Mick, "spend b1 b1 b1 b1"))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "spend b1 b1 b1"))
}
