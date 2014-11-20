package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCamouflage(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give steve card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Steve].Cards = []CardBase{&CardCamouflage{}}
	g.CurrentRoll = []int{}
	attack := 100
	for i := 0; i < attack; i++ {
		g.CurrentRoll = append(g.CurrentRoll, DieAttack)
	}
	g.Boards[Steve].Health = attack + 1
	cmd(t, g, Mick, "keep")
	assert.NotEqual(t, 1, g.Boards[Steve].Health)
}
