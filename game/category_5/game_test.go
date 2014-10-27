package category_5

import (
	"testing"

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

func TestSortCards(t *testing.T) {
	assert.Equal(t, []Card{1, 2, 3}, SortCards([]Card{3, 2, 1}))
}
