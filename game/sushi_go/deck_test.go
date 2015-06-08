package sushi_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeck(t *testing.T) {
	deck := Deck()
	assert.Len(t, deck, 108)
}

func TestSort(t *testing.T) {
	deck := []int{
		CardSquidNigiri,
		CardSashimi,
		CardMakiRoll1,
	}
	sorted := Sort(deck)
	assert.Equal(t, []int{
		CardSashimi,
		CardMakiRoll1,
		CardSquidNigiri,
	}, sorted)
	assert.Equal(t, []int{
		CardSquidNigiri,
		CardSashimi,
		CardMakiRoll1,
	}, deck)
}

func TestShuffle(t *testing.T) {
	deck := []int{
		CardSquidNigiri,
		CardSashimi,
		CardMakiRoll1,
	}
	shuffled := Shuffle(deck)
	assert.Len(t, shuffled, len(deck))
	for _, c := range deck {
		assert.Contains(t, shuffled, c)
	}
}
