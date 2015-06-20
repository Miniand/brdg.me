package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFireBreathingModifyAttackNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []int{FireBreathing}
	g.CurrentRoll = []int{}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 10, g.Boards[Steve].Health)
	assert.Equal(t, 10, g.Boards[BJ].Health)
	assert.Equal(t, 10, g.Boards[Walas].Health)
}

func TestCardFireBreathingModifyAttackWithAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []int{FireBreathing}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 7, g.Boards[Steve].Health)
	assert.Equal(t, 8, g.Boards[BJ].Health)
	assert.Equal(t, 7, g.Boards[Walas].Health)
}
