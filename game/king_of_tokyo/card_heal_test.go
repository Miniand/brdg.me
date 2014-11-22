package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHeal(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Health = 5
	g.Boards[Mick].Energy = 3
	g.Buyable = []CardBase{&CardHeal{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy heal")
	assert.Equal(t, 7, g.Boards[Mick].Health)
}
