package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardBackgroundDweller(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	g.Boards[0].Cards = []CardBase{&CardBackgroundDweller{}}
	g.CurrentRoll = []int{Die3}
	g.CheckRollComplete()
	assert.Equal(t, PhaseRoll, g.Phase)
	g.CurrentRoll = []int{Die2}
	g.CheckRollComplete()
	assert.Equal(t, PhaseBuy, g.Phase)
}
