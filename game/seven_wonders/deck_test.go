package seven_wonders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckAge1(t *testing.T) {
	d := DeckAge1()
	assert.Len(t, d, 49)
	for p := 3; p <= 7; p++ {
		assert.Len(t, FilterDeck(d, p), p*7)
	}
}

func TestDeckAge2(t *testing.T) {
	d := DeckAge2()
	assert.Len(t, d, 49)
	for p := 3; p <= 7; p++ {
		assert.Len(t, FilterDeck(d, p), p*7)
	}
}

func TestDeckAge3(t *testing.T) {
	d := DeckAge3()
	assert.Len(t, d, 40)
	for p := 3; p <= 7; p++ {
		assert.Len(t, FilterDeck(d, p), p*6-2)
	}
}
