package roll_through_the_ages

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_KeepSkulls_allDisasterSkip(t *testing.T) {
	g := &Game{}
	g.Start(helper.Players[:3])
	g.CurrentPlayer = helper.Mick
	g.RolledDice = []int{
		DiceSkull,
		DiceSkull,
		DiceSkull,
	}
	g.KeepSkulls()
	assert.Equal(t, helper.Mick, g.CurrentPlayer)
	assert.Equal(t, PhaseBuy, g.Phase)
}

func TestGame_KeepSkulls_allDisasterLeadership(t *testing.T) {
	g := &Game{}
	g.Start(helper.Players[:3])
	g.CurrentPlayer = helper.Mick
	g.Boards[helper.Mick].Developments[DevelopmentLeadership] = true
	g.RolledDice = []int{
		DiceSkull,
		DiceSkull,
		DiceSkull,
	}
	g.KeepSkulls()
	assert.Equal(t, helper.Mick, g.CurrentPlayer)
	assert.Equal(t, PhaseExtraRoll, g.Phase)
}
