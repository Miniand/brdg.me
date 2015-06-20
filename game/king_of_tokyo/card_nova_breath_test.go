package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardNovaBreath(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Boards[helper.Mick].Cards = []int{NovaBreath}
	g.CurrentRoll = []int{DieAttack, DieAttack}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	assert.Equal(t, 10, g.Boards[helper.Mick].Health)
	assert.Equal(t, 8, g.Boards[helper.Steve].Health)
	assert.Equal(t, 8, g.Boards[helper.BJ].Health)
}
