package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

func TestCardAlienMetabolismModifyAttackNoAttackDice(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve}))
	g.Boards[0].Cards = []CardBase{&CardAlienMetabolism{}}
	g.Boards[0].Energy = 5
	g.Buyable = []CardBase{&CardAlienMetabolism{}}
	g.Phase = PhaseBuy
	_, err := command.CallInCommands(Mick, g, "buy alien metabolism", g.Commands())
	assert.NoError(t, err)
	assert.Equal(t, 3, g.Boards[0].Energy)
}
