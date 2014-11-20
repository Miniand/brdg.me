package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCommuterTrain(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardCommuterTrain{}}
	g.Boards[Mick].Energy = 4
	g.Buyable = []CardBase{&CardCommuterTrain{}}
	g.Phase = PhaseBuy
	cmd(t, g, Mick, "buy commuter train")
	assert.Equal(t, 2, g.Boards[Mick].VP)
}
