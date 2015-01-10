package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardJets_Stay(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Tokyo[LocationTokyoCity] = Steve
	g.Boards[Steve].Cards = []CardBase{&CardJets{}}
	g.CurrentRoll = []int{
		DieAttack,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Steve, "stay")
	assert.Equal(t, 9, g.Boards[Steve].Health)
}

func TestCardJets_Leave(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Tokyo[LocationTokyoCity] = Steve
	g.Boards[Steve].Cards = []CardBase{&CardJets{}}
	g.CurrentRoll = []int{
		DieAttack,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Steve, "leave")
	assert.Equal(t, 10, g.Boards[Steve].Health)
}
