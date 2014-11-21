package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardEnergyHoarder(t *testing.T) {
	for _, cas := range []struct {
		energy, expectedVP int
	}{
		{0, 0},
		{1, 0},
		{5, 0},
		{6, 1},
		{11, 1},
		{12, 2},
	} {
		g := &Game{}
		assert.NoError(t, g.Start(names))
		g.Boards[Mick].Cards = []CardBase{&CardEnergyHoarder{}}
		g.Boards[Mick].Energy = cas.energy
		g.Phase = PhaseBuy
		cmd(t, g, Mick, "done")
		assert.Equal(t, cas.expectedVP, g.Boards[Mick].VP)
	}
}
