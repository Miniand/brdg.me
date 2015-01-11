package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardArmorPlatingModifyDamage1AttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give steve card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Steve].Cards = []CardBase{&CardArmorPlating{}}
	g.CurrentRoll = []int{
		DieAttack,
	}
	startingHealth := g.Boards[Steve].Health
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, startingHealth, g.Boards[Steve].Health)
}

func TestCardArmorPlatingModifyDamage2AttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Steve].Cards = []CardBase{&CardArmorPlating{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[Steve].Health
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, startingHealth-2, g.Boards[Steve].Health)
}
