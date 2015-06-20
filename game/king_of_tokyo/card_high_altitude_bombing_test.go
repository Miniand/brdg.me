package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHerbivore(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.FaceUpCards = []int{HighAltitudeBombing}
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.NoError(t, cmd(g, Mick, "buy high"))
	for p, _ := range g.Players {
		assert.Equal(t, 7, g.Boards[p].Health)
	}
}
