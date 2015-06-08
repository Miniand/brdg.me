package zombie_dice

import (
	"testing"

	h "github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_RenderForPlayer(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(h.Players[:2]))
	out, err := g.RenderForPlayer(h.Players[h.Mick])
	assert.NotEmpty(t, out)
	assert.NoError(t, err)
}
