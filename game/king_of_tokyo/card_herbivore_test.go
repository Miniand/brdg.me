package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHerbivore_VP(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardHerbivore{}}
	g.CurrentRoll = []int{
		Die1,
		Die2,
		Die3,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Mick, "done")
	assert.Equal(t, 1, g.Boards[Mick].VP)
}

func TestCardHerbivore_NoVP(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardHerbivore{}}
	g.Tokyo[LocationTokyoCity] = Mick
	g.CurrentRoll = []int{
		DieAttack,
		Die1,
		Die1,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Mick, "done")
	assert.Equal(t, 0, g.Boards[Mick].VP)
}
