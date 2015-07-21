package love_letter

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestCharPrince_Play_end(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	g.Hands[helper.Mick] = []int{Prince, Princess}
	g.Hands[helper.Steve] = []int{Prince}
	g.Protected[helper.Steve] = true
	g.Eliminated[helper.BJ] = true
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play prince mick"))
	assert.Equal(t, 2, g.Round)
	assert.Equal(t, 1, g.Points[helper.Steve])
	assert.Equal(t, helper.Steve, g.CurrentPlayer)
	assert.Len(t, g.Hands[helper.Mick], 1)
	assert.Len(t, g.Hands[helper.Steve], 2)
	assert.Len(t, g.Hands[helper.BJ], 1)
}
