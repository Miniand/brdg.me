package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardExtraHead(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	assert.Len(t, g.CurrentRoll, 6)
	g.Boards[Steve].Cards = []CardBase{&CardExtraHead{}}
	g.Phase = PhaseBuy
	assert.NoError(t, cmd(g, Mick, "done"))
	assert.Len(t, g.CurrentRoll, 7)
}
