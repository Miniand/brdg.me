package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFrenzy(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 7
	g.Buyable = []CardBase{&CardFrenzy{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy frenzy")
	assert.Equal(t, PhaseRoll, g.Phase)
	assert.Equal(t, Mick, g.CurrentPlayer)
}
