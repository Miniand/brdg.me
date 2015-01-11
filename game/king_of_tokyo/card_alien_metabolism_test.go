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
	g.FaceUpCards = []CardBase{
		&CardAlienMetabolism{},
		&CardAlienMetabolism{},
	}
	g.Phase = PhaseBuy
	// First purchase shouldn't be discounted
	assert.NoError(t, cmd(g, Mick, "buy alien metabolism"))
	assert.Equal(t, 5, g.Boards[Mick].Energy)
	// Second should
	assert.NoError(t, cmd(g, Mick, "buy alien metabolism"))
	assert.Equal(t, 3, g.Boards[Mick].Energy)
}
