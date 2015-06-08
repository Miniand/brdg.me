package sushi_go

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommand_Call(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.Equal(t, helper.Players[:2], g.WhoseTurn())

	// Mick plays a card
	mickCard := g.Hands[helper.Mick][0]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 1"))
	assert.Equal(t, []int{mickCard}, g.Playing[helper.Mick])
	// Mick hasn't played the dummy card yet so should still be their turn.
	assert.Equal(t, helper.Players[:2], g.WhoseTurn())

	// Steve plays a card
	steveCard := g.Hands[helper.Steve][8]
	assert.NoError(t, helper.Cmd(g, helper.Steve, "play 9"))
	assert.Equal(t, []string{helper.Players[helper.Mick]}, g.WhoseTurn())

	// Mick plays the dummy card
	dummyCard := g.Hands[helper.Mick][4]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "dummy 5"))
	// Plays should have happened now.
	assert.Equal(t, []int{mickCard}, g.Played[helper.Mick])
	assert.Equal(t, []int{steveCard}, g.Played[helper.Steve])
	assert.Equal(t, []int{dummyCard}, g.Played[Dummy])
	assert.Equal(t, helper.Players[:2], g.WhoseTurn())
}
