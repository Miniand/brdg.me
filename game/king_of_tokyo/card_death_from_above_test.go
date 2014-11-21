package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardDeathFromAboveOutsideTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 5
	g.Buyable = []CardBase{&CardDeathFromAbove{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy death from above")
	assert.Equal(t, 3, g.Boards[Mick].VP)
}

func TestCardDeathFromAboveInsideTokyo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Mick].Energy = 5
	g.Buyable = []CardBase{&CardDeathFromAbove{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy death from above")
	assert.Equal(t, 2, g.Boards[Mick].VP)
}
