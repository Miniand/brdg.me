package red7

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type handTest struct {
	input, expected []int
}

func TestHighestCard(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7"),
			expected: crds("r7"),
		},
	} {
		assert.Equal(t, c.expected, HighestCard(c.input))
	}
}

func TestCardsOfOneNumber(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7"),
			expected: crds("r7", "b7"),
		},
		{
			input:    crds("r6", "r7", "b6"),
			expected: crds("r6", "b6"),
		},
	} {
		assert.Equal(t, c.expected, CardsOfOneNumber(c.input))
	}
}

func TestCardsOfOneColor(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7"),
			expected: crds("r7", "r6"),
		},
		{
			input:    crds("r6", "b7", "b6", "r5"),
			expected: crds("b7", "b6"),
		},
	} {
		assert.Equal(t, c.expected, CardsOfOneColor(c.input))
	}
}

func TestMostEvenCards(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7"),
			expected: crds("r6"),
		},
		{
			input:    crds("r6", "b7", "b6", "r5"),
			expected: crds("r6", "b6"),
		},
	} {
		assert.Equal(t, c.expected, MostEvenCards(c.input))
	}
}

func TestCardsOfDifferentColors(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7", "y3", "y7"),
			expected: crds("r7", "y7", "b7"),
		},
	} {
		assert.Equal(t, c.expected, CardsOfDifferentColors(c.input))
	}
}

func TestCardsThatFormARun(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("r6", "r7", "b7", "y3", "y7"),
			expected: crds("r7", "r6"),
		},
		{
			input:    crds("r6", "b1", "r7", "g2", "b7", "y3", "y7"),
			expected: crds("y3", "g2", "b1"),
		},
	} {
		assert.Equal(t, c.expected, CardsThatFormARun(c.input))
	}
}

func TestMostCardsBelow4(t *testing.T) {
	for _, c := range []handTest{
		{
			input:    crds("b1", "b4", "r6", "g2", "r7", "b7", "y3", "y7"),
			expected: crds("y3", "g2", "b1"),
		},
	} {
		assert.Equal(t, c.expected, MostCardsBelow4(c.input))
	}
}

func TestLeader(t *testing.T) {
	leader, leaderPal := Leader([][]int{
		crds("y5", "b2"),
		crds("r5", "b2"),
		crds("g6"),
	})
	assert.Equal(t, leader, 1)
	assert.Equal(t, leaderPal, crds("r5", "b2"))
}
