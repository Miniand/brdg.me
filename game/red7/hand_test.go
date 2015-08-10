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
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
			},
			expected: []int{
				SuitRed | Rank7,
			},
		},
	} {
		assert.Equal(t, c.expected, HighestCard(c.input))
	}
}

func TestCardsOfOneNumber(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
			},
			expected: []int{
				SuitRed | Rank7,
				SuitBlue | Rank7,
			},
		},
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank6,
			},
			expected: []int{
				SuitRed | Rank6,
				SuitBlue | Rank6,
			},
		},
	} {
		assert.Equal(t, c.expected, CardsOfOneNumber(c.input))
	}
}

func TestCardsOfOneColor(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
			},
			expected: []int{
				SuitRed | Rank7,
				SuitRed | Rank6,
			},
		},
		{
			input: []int{
				SuitRed | Rank6,
				SuitBlue | Rank7,
				SuitBlue | Rank6,
				SuitRed | Rank5,
			},
			expected: []int{
				SuitBlue | Rank7,
				SuitBlue | Rank6,
			},
		},
	} {
		assert.Equal(t, c.expected, CardsOfOneColor(c.input))
	}
}

func TestMostEvenCards(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
			},
			expected: []int{
				SuitRed | Rank6,
			},
		},
		{
			input: []int{
				SuitRed | Rank6,
				SuitBlue | Rank7,
				SuitBlue | Rank6,
				SuitRed | Rank5,
			},
			expected: []int{
				SuitRed | Rank6,
				SuitBlue | Rank6,
			},
		},
	} {
		assert.Equal(t, c.expected, MostEvenCards(c.input))
	}
}

func TestCardsOfDifferentColors(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
				SuitYellow | Rank3,
				SuitYellow | Rank7,
			},
			expected: []int{
				SuitRed | Rank7,
				SuitYellow | Rank7,
				SuitBlue | Rank7,
			},
		},
	} {
		assert.Equal(t, c.expected, CardsOfDifferentColors(c.input))
	}
}

func TestCardsThatFormARun(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitRed | Rank6,
				SuitRed | Rank7,
				SuitBlue | Rank7,
				SuitYellow | Rank3,
				SuitYellow | Rank7,
			},
			expected: []int{
				SuitRed | Rank7,
				SuitRed | Rank6,
			},
		},
		{
			input: []int{
				SuitRed | Rank6,
				SuitBlue | Rank1,
				SuitRed | Rank7,
				SuitGreen | Rank2,
				SuitBlue | Rank7,
				SuitYellow | Rank3,
				SuitYellow | Rank7,
			},
			expected: []int{
				SuitYellow | Rank3,
				SuitGreen | Rank2,
				SuitBlue | Rank1,
			},
		},
	} {
		assert.Equal(t, c.expected, CardsThatFormARun(c.input))
	}
}

func TestMostCardsBelow4(t *testing.T) {
	for _, c := range []handTest{
		{
			input: []int{
				SuitBlue | Rank1,
				SuitBlue | Rank4,
				SuitRed | Rank6,
				SuitGreen | Rank2,
				SuitRed | Rank7,
				SuitBlue | Rank7,
				SuitYellow | Rank3,
				SuitYellow | Rank7,
			},
			expected: []int{
				SuitYellow | Rank3,
				SuitGreen | Rank2,
				SuitBlue | Rank1,
			},
		},
	} {
		assert.Equal(t, c.expected, MostCardsBelow4(c.input))
	}
}
