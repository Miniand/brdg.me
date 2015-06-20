package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCornerStore(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{}
	g.Boards[Mick].Energy = 3
	g.FaceUpCards = []int{CornerStore}
	g.Phase = PhaseBuy
	assert.NoError(t, cmd(g, Mick, "buy corner store"))
	assert.Equal(t, 1, g.Boards[Mick].VP)
}
