package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

func TestCardAlphaMonsterModifyAttackNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[0].Cards = []CardBase{&CardAlphaMonster{}}
	g.CurrentRoll = []int{}
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, 0, g.Boards[0].VP)
}

func TestCardAlphaMonsterModifyAttackWithAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = 0
	g.Boards[0].Cards = []CardBase{&CardAlphaMonster{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	_, err := command.CallInCommands(Mick, g, "keep", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, 1, g.Boards[0].VP)
}
