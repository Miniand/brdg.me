package jaipur

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Render(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	_, err := g.RenderForPlayer(g.Players[helper.Mick])
	assert.NoError(t, err)
}
