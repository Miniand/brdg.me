package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFreezeTimeWithout111(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardFreezeTime{}}
	g.CurrentRoll = []int{}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.NotEqual(t, Mick, g.CurrentPlayer)
}

func TestCardFreezeTimeWith111(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Cards = []CardBase{&CardFreezeTime{}}
	g.CurrentRoll = []int{
		Die1,
		Die1,
		Die1,
		Die2,
		Die2,
		Die3,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.Equal(t, Mick, g.CurrentPlayer)
	assert.Len(t, g.CurrentRoll, 5)
	g.CurrentRoll = []int{
		Die1,
		Die1,
		Die1,
		Die2,
		Die3,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.Equal(t, Mick, g.CurrentPlayer)
	assert.Len(t, g.CurrentRoll, 4)
	g.CurrentRoll = []int{}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.NotEqual(t, Mick, g.CurrentPlayer)
}
