package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardEvenBiggerPostPurchase(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 4
	startingHealth := g.Boards[Mick].Health
	g.Buyable = []CardBase{&CardEvenBigger{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy even bigger")
	assert.Equal(t, 12, g.Boards[Mick].MaxHealth())
	assert.Equal(t, startingHealth+2, g.Boards[Mick].Health)
}

func TestCardEvenBiggerHeal(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardEvenBigger{}}
	startingHealth := g.Boards[Mick].Health
	g.CurrentRoll = []int{
		DieHeal,
		DieHeal,
		DieHeal,
	}
	cmd(t, g, Mick, "keep")
	assert.Equal(t, startingHealth+2, g.Boards[Mick].Health)
}
