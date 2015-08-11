package red7

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCard(t *testing.T) {
	for _, c := range []struct {
		input           string
		expectedSuccess bool
		expectedCard    int
	}{
		{"b5", true, SuitBlue | Rank5},
		{"a5", false, 0},
		{"r0", false, 0},
		{"red5", false, 0},
		{"r8", false, 0},
		{"r10", false, 0},
		{"o4", true, SuitOrange | Rank4},
	} {
		card, ok := ParseCard(c.input)
		assert.Equal(t, c.expectedSuccess, ok)
		if c.expectedSuccess {
			assert.Equal(t, c.expectedCard, card)
		}
	}
}
