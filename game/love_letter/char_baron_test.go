package love_letter

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCharBaron_Play_win(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Hands[helper.Mick] = []int{Baron, King}
	g.Hands[helper.Steve] = []int{Prince}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play baron steve"))
	assert.Equal(t, []int{King}, g.Hands[helper.Mick])
	assert.False(t, g.Eliminated[helper.Mick])
	assert.True(t, g.Eliminated[helper.Steve])
}

func TestCharBaron_Play_tie(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Hands[helper.Mick] = []int{Baron, Prince}
	g.Hands[helper.Steve] = []int{Prince}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play baron steve"))
	assert.Equal(t, []int{Prince}, g.Hands[helper.Mick])
	assert.False(t, g.Eliminated[helper.Mick])
	assert.False(t, g.Eliminated[helper.Steve])
}

func TestCharBaron_Play_lose(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Hands[helper.Mick] = []int{Baron, Prince}
	g.Hands[helper.Steve] = []int{King}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play baron steve"))
	assert.Equal(t, []int{}, g.Hands[helper.Mick])
	assert.True(t, g.Eliminated[helper.Mick])
	assert.False(t, g.Eliminated[helper.Steve])
}

func TestCharBaron_Play_double(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Hands[helper.Mick] = []int{Baron, Baron}
	g.Hands[helper.Steve] = []int{Guard}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play baron steve"))
	assert.Equal(t, []int{Baron}, g.Hands[helper.Mick])
	assert.False(t, g.Eliminated[helper.Mick])
	assert.True(t, g.Eliminated[helper.Steve])
}
