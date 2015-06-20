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
	g.FaceUpCards = []int{EvenBigger}
	g.Phase = PhaseBuy
	assert.NoError(t, cmd(g, Mick, "buy even bigger"))
	assert.Equal(t, 12, g.Boards[Mick].MaxHealth())
	assert.Equal(t, startingHealth+2, g.Boards[Mick].Health)
}

func TestCardEvenBiggerHeal(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{EvenBigger}
	startingHealth := g.Boards[Mick].Health
	g.CurrentRoll = []int{
		DieHeal,
		DieHeal,
		DieHeal,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, startingHealth+2, g.Boards[Mick].Health)
}
