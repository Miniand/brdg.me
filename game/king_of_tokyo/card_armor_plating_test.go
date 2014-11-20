package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

func TestCardArmorPlatingModifyDamage1AttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[1].Cards = []CardBase{&CardArmorPlating{}}
	g.CurrentRoll = []int{
		DieAttack,
	}
	startingHealth := g.Boards[1].Health
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, startingHealth, g.Boards[1].Health)
}

func TestCardArmorPlatingModifyDamage2AttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[1].Cards = []CardBase{&CardArmorPlating{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[1].Health
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, startingHealth-2, g.Boards[1].Health)
}
