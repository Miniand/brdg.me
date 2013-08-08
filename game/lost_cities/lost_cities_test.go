package lost_cities

import (
	"github.com/beefsack/brdg.me/game/card"
	"testing"
)

// Build a game by hand for testing purposes.  Each player has a full hand, half
// of the discard.SuitRankCard stacks have card.SuitRankCards, and there are two card.SuitRankCards in the draw pile.
func mockGame(t *testing.T) *Game {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Fatal(err)
	}
	game.CurrentlyMoving = 0
	// Set Mick's hand
	game.Board.PlayerHands[0] = card.Deck{
		card.SuitRankCard{
			Suit: SUIT_BLUE,
			Rank: 6,
		},
		card.SuitRankCard{
			Suit: SUIT_BLUE,
			Rank: 8,
		},
		card.SuitRankCard{
			Suit: SUIT_RED,
			Rank: 4,
		},
		card.SuitRankCard{
			Suit: SUIT_RED,
			Rank: 5,
		},
		card.SuitRankCard{
			Suit: SUIT_YELLOW,
			Rank: 0,
		},
		card.SuitRankCard{
			Suit: SUIT_YELLOW,
			Rank: 3,
		},
		card.SuitRankCard{
			Suit: SUIT_GREEN,
			Rank: 2,
		},
		card.SuitRankCard{
			Suit: SUIT_WHITE,
			Rank: 10,
		},
	}
	// Set Steve's hand
	game.Board.PlayerHands[1] = card.Deck{
		card.SuitRankCard{
			Suit: SUIT_BLUE,
			Rank: 7,
		},
		card.SuitRankCard{
			Suit: SUIT_BLUE,
			Rank: 9,
		},
		card.SuitRankCard{
			Suit: SUIT_RED,
			Rank: 0,
		},
		card.SuitRankCard{
			Suit: SUIT_RED,
			Rank: 10,
		},
		card.SuitRankCard{
			Suit: SUIT_YELLOW,
			Rank: 4,
		},
		card.SuitRankCard{
			Suit: SUIT_YELLOW,
			Rank: 7,
		},
		card.SuitRankCard{
			Suit: SUIT_GREEN,
			Rank: 4,
		},
		card.SuitRankCard{
			Suit: SUIT_WHITE,
			Rank: 8,
		},
	}
	// Just set the draw pile to have a couple of card.SuitRankCards so we can finish the
	// round quickly for testing.
	game.Board.DrawPile = card.Deck{
		card.SuitRankCard{
			Suit: SUIT_WHITE,
			Rank: 0,
		},
		card.SuitRankCard{
			Suit: SUIT_WHITE,
			Rank: 0,
		},
		card.SuitRankCard{
			Suit: SUIT_WHITE,
			Rank: 0,
		},
	}
	return game
}

func TestPlayFullGame(t *testing.T) {
	game := mockGame(t)
	if game.IsEndOfRound() || game.IsFinished() {
		t.Fatal("Why is it the end of the round if we've just started?")
	}
	if game.TurnPhase != TURN_PHASE_PLAY_OR_DISCARD {
		t.Fatal("The turn phase isn't for the player to play or discard.SuitRankCard")
	}
	// Mick discard.SuitRankCards red 5
	err := game.PlayerAction("Mick", "discard.SuitRankCard", []string{"r5"})
	if err != nil {
		t.Fatal(err)
	}
	// Let's check to make sure it actually happened
	if len(game.Board.PlayerHands[0]) != 7 {
		t.Fatal("Red 5 wasn't removed from Mick's hand when he discard.SuitRankCarded")
	}
	if len(game.Board.DiscardPiles[SUIT_RED]) != 1 {
		t.Fatal("The red discard.SuitRankCard pile doesn't have any card.SuitRankCards in it")
	}
	firstRed := game.Board.DiscardPiles[SUIT_RED][0].(card.SuitRankCard)
	if firstRed.Suit != SUIT_RED && firstRed.Rank != 5 {
		t.Fatal("Red 5 wasn't discard.SuitRankCarded onto the red discard.SuitRankCard pile")
	}
	if game.TurnPhase != TURN_PHASE_DRAW {
		t.Fatal("The turn phase didn't change to DRAW")
	}
	// Steve tries to butt in, but he shouldn't be allowed cos it's not his
	// turn!
	err = game.PlayerAction("Steve", "draw", []string{})
	if err == nil {
		t.Fatal(
			"Steve was allowed to draw a card.SuitRankCard even though it wasn't his turn!")
	}
	// More to come!
}
