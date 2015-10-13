package seven_wonders_duel

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
}

func TestOpponent(t *testing.T) {
	assert.Equal(t, 1, Opponent(0))
	assert.Equal(t, 0, Opponent(1))
}
