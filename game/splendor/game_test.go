package splendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func cmd(g *Game, player int, input string) error {
	return nil
}

func TestParseLoc(t *testing.T) {
	row, col, err := ParseLoc("1A")
	assert.NoError(t, err)
	assert.Equal(t, 0, row)
	assert.Equal(t, 0, col)

	row, col, err = ParseLoc("4i")
	assert.NoError(t, err)
	assert.Equal(t, 3, row)
	assert.Equal(t, 8, col)

	row, col, err = ParseLoc("0B")
	assert.Error(t, err)

	row, col, err = ParseLoc("5B")
	assert.Error(t, err)
}
