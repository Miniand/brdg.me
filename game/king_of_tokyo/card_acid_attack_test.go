package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

func TestCardAcidAttackModifyAttackNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[0].Cards = []CardBase{&CardAcidAttack{}}
	g.CurrentRoll = []int{}
	startingHealth := g.Boards[1].Health
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, startingHealth-1, g.Boards[1].Health)
}

func TestCardAcidAttackModifyAttackWithAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[0].Cards = []CardBase{&CardAcidAttack{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[1].Health
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, startingHealth-3, g.Boards[1].Health)
}
