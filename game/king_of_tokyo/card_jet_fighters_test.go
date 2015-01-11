package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardJetFighters(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.FaceUpCards = []CardBase{&CardJetFighters{}}
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "buy jet"))
	assert.Equal(t, 6, g.Boards[Mick].Health)
	assert.Equal(t, 5, g.Boards[Mick].VP)
}
