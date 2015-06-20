package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardDedicatedNewsTeamOutsideTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 6
	g.Boards[Mick].Cards = []int{DedicatedNewsTeam}
	g.FaceUpCards = []int{AcidAttack}
	g.Phase = PhaseBuy
	assert.NoError(t, cmd(g, Mick, "buy acid attack"))
	assert.Equal(t, 1, g.Boards[Mick].VP)
}
