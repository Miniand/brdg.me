package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardEaterOfTheDead(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	// Put Mick in Tokyo and give card
	g.Tokyo[LocationTokyoCity] = Mick
	g.Boards[Steve].Cards = []int{EaterOfTheDead}
	g.Boards[Steve].Health = 3
	g.Boards[BJ].Cards = []int{EaterOfTheDead}
	g.CurrentRoll = []int{
		DieAttack,
		DieAttack,
		DieAttack,
	}
	assert.NoError(t, cmd(g, Mick, "keep"))
	assert.Equal(t, 3, g.Boards[Steve].VP)
	assert.Equal(t, 3, g.Boards[BJ].VP)
}
