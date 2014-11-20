package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

func TestCardApartmentBuilding(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	g.Boards[0].Cards = []CardBase{&CardApartmentBuilding{}}
	g.Boards[0].Energy = 5
	g.Buyable = []CardBase{&CardApartmentBuilding{}}
	g.Phase = PhaseBuy
	_, err := command.CallInCommands(Mick, g, "buy apartment building", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, 3, g.Boards[0].VP)
}
