package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardGourmet(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardGourmet{}}
	g.CurrentRoll = []int{
		Die1,
		Die1,
		Die1,
	}
	cmd(t, g, Mick, "keep")
	assert.Equal(t, 3, g.Boards[Mick].VP)
}
