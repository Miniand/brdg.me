package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardEnergize(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{}
	g.Boards[Mick].Energy = 8
	g.Buyable = []CardBase{&CardEnergize{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy energize")
	assert.Equal(t, 9, g.Boards[Mick].Energy)
}
