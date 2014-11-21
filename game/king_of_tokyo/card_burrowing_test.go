package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardBurrowingModifyAttackOutsideTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Steve in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Steve
	g.Boards[Mick].Cards = []CardBase{&CardBurrowing{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[Steve].Health
	cmd(t, g, Mick, "keep")
	cmd(t, g, Steve, "stay")
	assert.Equal(t, startingHealth-2, g.Boards[Steve].Health)
}

func TestCardBurrowingModifyAttackInTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardBurrowing{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[Steve].Health
	cmd(t, g, Mick, "keep")
	assert.Equal(t, startingHealth-3, g.Boards[Steve].Health)
}

func TestCardBurrowingDamageWhenLeavingTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Steve in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Steve
	g.Boards[Steve].Cards = []CardBase{&CardBurrowing{}}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
	}
	startingHealth := g.Boards[Mick].Health
	cmd(t, g, Mick, "keep")
	cmd(t, g, Steve, "leave")
	assert.Equal(t, startingHealth-1, g.Boards[Mick].Health)
}
