package age_of_war

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_RenderForPlayer(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	_, err := g.RenderForPlayer(helper.Players[helper.Mick])
	assert.NoError(t, err)
}
