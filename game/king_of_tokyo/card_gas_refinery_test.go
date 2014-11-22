package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGasRefinery(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 6
	g.Buyable = []CardBase{&CardGasRefinery{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy gas refinery")
	assert.Equal(t, 2, g.Boards[Mick].VP)
	assert.Equal(t, 10, g.Boards[Mick].Health)
	assert.Equal(t, 7, g.Boards[Steve].Health)
	assert.Equal(t, 7, g.Boards[BJ].Health)
	assert.Equal(t, 7, g.Boards[Walas].Health)
}
