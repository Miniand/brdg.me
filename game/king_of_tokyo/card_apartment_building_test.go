package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardApartmentBuilding(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{}
	g.Boards[Mick].Energy = 5
	g.FaceUpCards = []int{ApartmentBuilding}
	g.Phase = PhaseBuy
	assert.NoError(t, cmd(g, Mick, "buy apartment building"))
	assert.Equal(t, 3, g.Boards[Mick].VP)
}
