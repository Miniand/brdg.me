package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFriendOfChildren(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []int{FriendOfChildren}
	g.CurrentRoll = []int{
		DieEnergy,
		DieEnergy,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 3, g.Boards[Mick].Energy)
}
