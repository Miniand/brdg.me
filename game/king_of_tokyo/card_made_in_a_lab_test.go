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
	cmd(t, g, Mick, "buy "+g.Deck[0].Name())
}
