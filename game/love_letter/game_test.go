package love_letter

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.NoError(t, g.Start(helper.Players[:3]))
	assert.NoError(t, g.Start(helper.Players[:4]))
}

func TestGame_Decode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	data, err := g.Encode()
	assert.NoError(t, err)
	assert.NoError(t, g.Decode(data))
}

func TestGame_RenderForPlayer(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	for _, p := range g.Players {
		_, err := g.RenderForPlayer(p)
		assert.NoError(t, err)
	}
}

func TestGame_IsFinished(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[2] - 1
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[2]
	assert.True(t, g.IsFinished())

	g = &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[3] - 1
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[3]
	assert.True(t, g.IsFinished())

	g = &Game{}
	assert.NoError(t, g.Start(helper.Players[:4]))
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[4] - 1
	assert.False(t, g.IsFinished())
	g.Points[helper.Mick] = endScores[4]
	assert.True(t, g.IsFinished())
}

func TestGame_Winners(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	assert.Equal(t, []string{}, g.Winners())
	g.Points[helper.Steve] = 10
	assert.Equal(t, []string{helper.Players[helper.Steve]}, g.Winners())
}
