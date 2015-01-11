package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardMadeInALab(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardMadeInALab{}}
	g.Boards[Mick].Energy = 20
	cmd(t, g, Mick, "keep")

	extraCard := g.Deck[0]
	assert.Contains(t, g.Buyable(Mick), extraCard)
	assert.NotContains(t, g.Buyable(Steve), extraCard)
	cmd(t, g, Mick, "buy "+extraCard.Name())
}
