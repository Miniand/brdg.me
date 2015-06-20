package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardBackgroundDweller(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{BackgroundDweller}
	g.CurrentRoll = []int{Die3}
	g.CheckRollComplete()
	assert.Equal(t, PhaseRoll, g.Phase)
	g.CurrentRoll = []int{Die2}
	g.CheckRollComplete()
	assert.Equal(t, PhaseBuy, g.Phase)
}
