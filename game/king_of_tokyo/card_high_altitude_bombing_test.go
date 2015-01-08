package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHerbivore(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Buyable = []CardBase{&CardHighAltitudeBombing{}}
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Mick, "buy high")
	for p, _ := range g.Players {
		assert.Equal(t, 7, g.Boards[p].Health)
	}
}
