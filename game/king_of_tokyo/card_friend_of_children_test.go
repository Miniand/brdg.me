package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFriendOfChildren(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardFriendOfChildren{}}
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
	}
	cmd(t, g, Mick, "keep")
	assert.Equal(t, 3, g.Boards[Mick].Energy)
}
