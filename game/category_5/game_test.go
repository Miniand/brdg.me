package category_5

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_DrawCards(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{"Mick", "Steve"}))
	g.Discard = g.DrawCards(75)
	assert.Len(t, g.Discard, 75)
	assert.Len(t, g.Deck, 5)
	assert.Len(t, g.DrawCards(10), 10)
	assert.Len(t, g.Discard, 0)
	assert.Len(t, g.Deck, 70)
}

func TestAutoPlayLastCard(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	g.Board = [4][]Card{
		{1},
		{2},
		{3},
		{4},
	}
	g.Hands = map[int][]Card{
		helper.Mick:  {5, 6},
		helper.Steve: {7, 8},
	}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 5"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "play 7"))
	assert.Len(t, g.Hands[helper.Mick], 10)
}

func TestSortCards(t *testing.T) {
	assert.Equal(t, []Card{1, 2, 3}, SortCards([]Card{3, 2, 1}))
}
