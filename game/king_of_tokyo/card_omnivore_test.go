package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCardOmnivore(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.CurrentRoll = []int{Die1, Die2, Die3, Die3}
	g.Boards[helper.Mick].Cards = []int{Omnivore}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "keep"))
	assert.Equal(t, 2, g.Boards[helper.Mick].VP)
}
