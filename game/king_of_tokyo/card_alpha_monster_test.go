package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardAlphaMonsterNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardAlphaMonster{}}
	g.CurrentRoll = []int{}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 0, g.Boards[Mick].VP)
}

func TestCardAlphaMonsterWithAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardAlphaMonster{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 1, g.Boards[Mick].VP)
}
