package sushi_go

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestPlayCommand_Call(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	assert.Equal(t, helper.Players[:3], g.WhoseTurn())

	// Mick plays a card
	mickHandLen := len(g.Hands[helper.Mick])
	mickCard := g.Hands[helper.Mick][0]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 1"))
	assert.Len(t, g.Hands[helper.Mick], mickHandLen-1)
	assert.Equal(t, []int{mickCard}, g.Playing[helper.Mick])
	assert.Equal(t, helper.Players[1:3], g.WhoseTurn())

	// BJ plays a card
	bjHandLen := len(g.Hands[helper.BJ])
	bjCard := g.Hands[helper.BJ][1]
	assert.NoError(t, helper.Cmd(g, helper.BJ, "play 2"))
	assert.Len(t, g.Hands[helper.BJ], bjHandLen-1)
	assert.Equal(t, []int{bjCard}, g.Playing[helper.BJ])
	assert.Equal(t, helper.Players[1:2], g.WhoseTurn())

	// Steve plays a card
	steveHandLen := len(g.Hands[helper.Steve])
	steveCard := g.Hands[helper.Steve][8]
	assert.NoError(t, helper.Cmd(g, helper.Steve, "play 9"))
	assert.Len(t, g.Hands[helper.Steve], steveHandLen-1)
	// Plays should have happened now.
	assert.Equal(t, []int{mickCard}, g.Played[helper.Mick])
	assert.Equal(t, []int{bjCard}, g.Played[helper.BJ])
	assert.Equal(t, []int{steveCard}, g.Played[helper.Steve])
	assert.Equal(t, helper.Players[:3], g.WhoseTurn())
}
