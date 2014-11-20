package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardAlienMetabolism(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{}
	g.Boards[Mick].Energy = 8
	g.Buyable = []CardBase{
		&CardAlienMetabolism{},
		&CardAlienMetabolism{},
	}
	g.Phase = PhaseBuy
	// First purchase shouldn't be discounted
	cmd(t, g, Mick, "buy alien metabolism")
	assert.Equal(t, 5, g.Boards[Mick].Energy)
	// Second should
	cmd(t, g, Mick, "buy alien metabolism")
	assert.Equal(t, 3, g.Boards[Mick].Energy)
}
