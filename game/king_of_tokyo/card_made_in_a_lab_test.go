package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buyableCardsToCards(bc []BuyableCard) []CardBase {
	cards := []CardBase{}
	for _, c := range bc {
		cards = append(cards, c.Card)
	}
	return cards
}

func TestCardMadeInALab(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardMadeInALab{}}
	g.Boards[Mick].Energy = 20
	assert.NoError(t, cmd(g, Mick, "keep"))
	extraCard := g.Deck[0]
	assert.NoError(t, cmd(g, Mick, "buy "+extraCard.Name()))
	assert.False(t, extraCard == g.Deck[0])
}
