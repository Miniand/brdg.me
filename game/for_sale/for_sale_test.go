package for_sale

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

const (
	Mick  = "Mick"
	Steve = "Steve"
	BJ    = "BJ"
)

func cmd(g *Game, player, input string) error {
	_, err := command.CallInCommands(player, g, input, g.Commands())
	return err
}

func TestFullGame(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start([]string{Mick, Steve, BJ}))
	// Set the state of the game to sorted decks
	_, g.BuildingDeck = BuildingDeck().PopN(2)
	_, g.ChequeDeck = ChequeDeck().PopN(2)
	g.OpenCards, g.BuildingDeck = g.BuildingDeck.PopN(3)
	// Play a round of buying
	assert.Equal(t, []string{Mick}, g.WhoseTurn())
	assert.NoError(t, cmd(g, Mick, "bid 3"))
	assert.Equal(t, []string{Steve}, g.WhoseTurn())
	assert.Error(t, cmd(g, Steve, "bid 3"))
	assert.NoError(t, cmd(g, Steve, "bid 4"))
	assert.Equal(t, []string{BJ}, g.WhoseTurn())
	assert.NoError(t, cmd(g, BJ, "pass"))
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 17},
		card.SuitRankCard{Rank: 18},
	}, g.OpenCards)
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 16},
	}, g.Hands[2])
	assert.Equal(t, 15, g.Chips[2])
	assert.Equal(t, []string{Mick}, g.WhoseTurn())
	assert.NoError(t, cmd(g, Mick, "pass"))
	assert.Equal(t, 14, g.Chips[0])
	assert.Equal(t, 11, g.Chips[1])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 17},
	}, g.Hands[0])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 18},
	}, g.Hands[1])
	assert.Equal(t, []string{Steve}, g.WhoseTurn())
	// One more buying phase so each player has 2 buildings.
	assert.NoError(t, cmd(g, Steve, "pass"))
	assert.NoError(t, cmd(g, BJ, "pass"))
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 15},
		card.SuitRankCard{Rank: 17},
	}, g.Hands[0])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 13},
		card.SuitRankCard{Rank: 18},
	}, g.Hands[1])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 14},
		card.SuitRankCard{Rank: 16},
	}, g.Hands[2])
	// End the buying phase early and shorten the selling phase.
	g.BuildingDeck = card.Deck{}
	_, g.ChequeDeck = g.ChequeDeck.PopN(12)
	g.OpenCards = card.Deck{}
	g.StartRound()
	assert.Equal(t, []string{Mick, Steve, BJ}, g.WhoseTurn())
	// Play a round of selling
	assert.Error(t, cmd(g, BJ, "play 18"))
	assert.NoError(t, cmd(g, BJ, "play 16"))
	assert.Equal(t, []string{Mick, Steve}, g.WhoseTurn())
	assert.NoError(t, cmd(g, Steve, "play 18"))
	assert.Equal(t, []string{Mick}, g.WhoseTurn())
	assert.NoError(t, cmd(g, Mick, "play 17"))
	// Because there were only two cards each, assume that the last cards were
	// automatically played.
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 5},
		card.SuitRankCard{Rank: 3},
	}, g.Cheques[0])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 6},
		card.SuitRankCard{Rank: 0},
	}, g.Cheques[1])
	assert.Equal(t, card.Deck{
		card.SuitRankCard{Rank: 4},
		card.SuitRankCard{Rank: 0},
	}, g.Cheques[2])
	// Check the game ended
	assert.True(t, g.IsFinished())
	assert.Equal(t, []string{Mick}, g.Winners())
}
