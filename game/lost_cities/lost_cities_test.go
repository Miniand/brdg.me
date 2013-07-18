package lost_cities

import (
	"testing"
)

// Build a game by hand for testing purposes.  Each player has a full hand, half
// of the discard stacks have cards, and there are two cards in the draw pile.
func mockGame() *Game {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		panic("Unable to start the game: " + err.Error())
	}
	return game
}

func TestStart(t *testing.T) {
	game := mockGame()
	if game.PlayerList()[0] != "Mick" {
		t.Error("Mick is not a player")
	}
	if game.PlayerList()[1] != "Steve" {
		t.Error("Steve is not a player")
	}
	if len(game.Board.PlayerHands[0].CardsBySuit[SUIT_RED])+
		len(game.Board.PlayerHands[0].CardsBySuit[SUIT_GREEN])+
		len(game.Board.PlayerHands[0].CardsBySuit[SUIT_BLUE])+
		len(game.Board.PlayerHands[0].CardsBySuit[SUIT_WHITE])+
		len(game.Board.PlayerHands[0].CardsBySuit[SUIT_YELLOW]) != 8 {
		t.Error("Mick's hand isn't 8 cards in total")
	}
	if len(game.Board.PlayerHands[1].CardsBySuit[SUIT_RED])+
		len(game.Board.PlayerHands[1].CardsBySuit[SUIT_GREEN])+
		len(game.Board.PlayerHands[1].CardsBySuit[SUIT_BLUE])+
		len(game.Board.PlayerHands[1].CardsBySuit[SUIT_WHITE])+
		len(game.Board.PlayerHands[1].CardsBySuit[SUIT_YELLOW]) != 8 {
		t.Error("Steve's hand isn't 8 cards in total")
	}
	if game.IsFinished() {
		t.Error("Game is already finished but we just started!")
	}
}

func TestAllCards(t *testing.T) {
	game := mockGame()
	allCards := game.AllCards()
	if len(allCards) != 60 {
		t.Error("There aren't 60 cards in the full deck")
	}
}
