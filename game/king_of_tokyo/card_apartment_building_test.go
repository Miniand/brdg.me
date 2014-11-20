package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardApartmentBuilding(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardApartmentBuilding{}}
	g.Boards[Mick].Energy = 5
	g.Buyable = []CardBase{&CardApartmentBuilding{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy apartment building")
	assert.Equal(t, 3, g.Boards[Mick].VP)
}
