package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardEvacuationOrders(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{}
	g.Boards[Mick].Energy = 7
	g.Boards[Mick].VP = 5
	g.Boards[Steve].VP = 7
	g.Boards[BJ].VP = 3
	g.Buyable = []CardBase{&CardEvacuationOrders{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy evacuation orders")
	assert.Equal(t, 5, g.Boards[Mick].VP)
	assert.Equal(t, 2, g.Boards[Steve].VP)
	assert.Equal(t, 0, g.Boards[BJ].VP)
}
