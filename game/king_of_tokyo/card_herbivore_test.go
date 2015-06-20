package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHerbivore_VP(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{Herbivore}
	g.CurrentRoll = []int{
		Die1,
		Die2,
		Die3,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.Equal(t, 1, g.Boards[Mick].VP)
}

func TestCardHerbivore_NoVP(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{Herbivore}
	g.Tokyo[LocationTokyoCity] = Mick
	g.CurrentRoll = []int{
		DieAttack,
		Die1,
		Die1,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.Equal(t, 0, g.Boards[Mick].VP)
}
