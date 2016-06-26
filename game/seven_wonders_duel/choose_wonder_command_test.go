package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_CanChooseWonder(t *testing.T) {
	g := Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	for _, tc := range []struct {
		numRemaining, player int
	}{
		{8, 0},
		{7, 1},
		{6, 1},
		{5, 0},
		{4, 1},
		{3, 0},
		{2, 0},
		{1, 1},
	} {
		g.RemainingWonders = make([]int, tc.numRemaining)
		assert.True(t, g.CanChooseWonder(tc.player), fmt.Sprintf("%#v", tc))
		assert.False(t, g.CanChooseWonder(Opponent(tc.player)), fmt.Sprintf("%#v", tc))
	}
}
