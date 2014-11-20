package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCompleteDestructionIncorrectDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardCompleteDestruction{}}
	g.CurrentRoll = []int{
		Die1,
		Die2,
		Die3,
		DieEnergy,
		DieHeal,
		DieHeal,
	}
	cmd(t, g, Mick, "keep")
	assert.Equal(t, 0, g.Boards[Mick].VP)
}

func TestCardCompleteDestructionCorrectDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardCompleteDestruction{}}
	g.CurrentRoll = []int{
		Die1,
		Die2,
		Die3,
		DieEnergy,
		DieHeal,
		DieAttack,
	}
	cmd(t, g, Mick, "keep")
	assert.Equal(t, 9, g.Boards[Mick].VP)
}
