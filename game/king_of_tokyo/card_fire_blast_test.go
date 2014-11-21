package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFireBlast(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Energy = 3
	g.Buyable = []CardBase{&CardFireBlast{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy fire blast")
	assert.Equal(t, 10, g.Boards[Mick].Health)
	assert.Equal(t, 8, g.Boards[Steve].Health)
	assert.Equal(t, 8, g.Boards[BJ].Health)
}
