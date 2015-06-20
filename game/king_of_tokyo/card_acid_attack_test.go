package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardAcidAttackModifyAttackNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []int{AcidAttack}
	g.CurrentRoll = []int{}
	startingHealth := g.Boards[Steve].Health
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, startingHealth-1, g.Boards[Steve].Health)
}

func TestCardAcidAttackModifyAttackWithAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []int{AcidAttack}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[Steve].Health
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, startingHealth-3, g.Boards[Steve].Health)
}
