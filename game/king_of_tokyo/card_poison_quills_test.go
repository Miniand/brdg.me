package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardPoisonQuills(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Tokyo[LocationTokyoCity] = helper.Steve
	g.Boards[helper.Mick].Cards = []int{PoisonQuills}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
		DieAttack,
		Die2,
		Die2,
		Die2,
	}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	assert.NoError(t, helper.Cmd(g, helper.Steve, "leave"))
	assert.Equal(t, 10, g.Boards[helper.Mick].Health)
	assert.Equal(t, 5, g.Boards[helper.Steve].Health)
	assert.Equal(t, 10, g.Boards[helper.BJ].Health)
}
