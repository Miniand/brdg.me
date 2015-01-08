package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardItHasAChild(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Buyable = []CardBase{&CardItHasAChild{}}
	g.Boards[Mick].Energy = 5
	g.Boards[Mick].VP = 5
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieEnergy,
		DieAttack,
	}
	cmd(t, g, Mick, "keep")
	cmd(t, g, Mick, "buy it has")
	cmd(t, g, Mick, "done")
	assert.Equal(t, Mick, g.Tokyo[LocationTokyoCity])
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
		DieAttack,
	}
	cmd(t, g, Steve, "keep")
	cmd(t, g, Mick, "leave")
	assert.Equal(t, 10, g.Boards[Mick].Health)
	assert.Len(t, g.Boards[Mick].Cards, 0)
	assert.Equal(t, 5, g.Boards[Mick].Energy)
	assert.Equal(t, 0, g.Boards[Mick].VP)
}
