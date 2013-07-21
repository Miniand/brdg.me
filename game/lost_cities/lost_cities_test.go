package lost_cities

import (
	"testing"
)

// Build a game by hand for testing purposes.  Each player has a full hand, half
// of the discard stacks have cards, and there are two cards in the draw pile.
func mockGame(t *testing.T) *Game {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Fatal(err)
	}
	game.CurrentlyMoving = 0
	// Set Mick's hand
	game.Board.PlayerHands[0] = []Card{
		Card{
			Suit:  SUIT_BLUE,
			Value: 6,
		},
		Card{
			Suit:  SUIT_BLUE,
			Value: 8,
		},
		Card{
			Suit:  SUIT_RED,
			Value: 4,
		},
		Card{
			Suit:  SUIT_RED,
			Value: 5,
		},
		Card{
			Suit:  SUIT_YELLOW,
			Value: 0,
		},
		Card{
			Suit:  SUIT_YELLOW,
			Value: 3,
		},
		Card{
			Suit:  SUIT_GREEN,
			Value: 2,
		},
		Card{
			Suit:  SUIT_WHITE,
			Value: 10,
		},
	}
	// Set Steve's hand
	game.Board.PlayerHands[1] = []Card{
		Card{
			Suit:  SUIT_BLUE,
			Value: 7,
		},
		Card{
			Suit:  SUIT_BLUE,
			Value: 9,
		},
		Card{
			Suit:  SUIT_RED,
			Value: 0,
		},
		Card{
			Suit:  SUIT_RED,
			Value: 10,
		},
		Card{
			Suit:  SUIT_YELLOW,
			Value: 4,
		},
		Card{
			Suit:  SUIT_YELLOW,
			Value: 7,
		},
		Card{
			Suit:  SUIT_GREEN,
			Value: 4,
		},
		Card{
			Suit:  SUIT_WHITE,
			Value: 8,
		},
	}
	// Just set the draw pile to have a couple of cards so we can finish the
	// round quickly for testing.
	game.Board.DrawPile = []Card{
		Card{
			Suit:  SUIT_WHITE,
			Value: 0,
		},
		Card{
			Suit:  SUIT_WHITE,
			Value: 0,
		},
		Card{
			Suit:  SUIT_WHITE,
			Value: 0,
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
		t.Fatal("The turn phase isn't for the player to play or discard")
	}
	// Mick discards red 5
	err := game.PlayerAction("Mick", "discard", []string{"r5"})
	if err != nil {
		t.Fatal(err)
	}
	// Let's check to make sure it actually happened
	if len(game.Board.PlayerHands[0]) != 7 {
		t.Fatal("Red 5 wasn't removed from Mick's hand when he discarded")
	}
	if len(game.Board.DiscardPiles[SUIT_RED]) != 1 {
		t.Fatal("The red discard pile doesn't have any cards in it")
	}
	if game.Board.DiscardPiles[SUIT_RED][0].Suit != SUIT_RED &&
		game.Board.DiscardPiles[SUIT_RED][0].Value != 5 {
		t.Fatal("Red 5 wasn't discarded onto the red discard pile")
	}
	if game.TurnPhase != TURN_PHASE_DRAW {
		t.Fatal("The turn phase didn't change to DRAW")
	}
	// Steve tries to butt in, but he shouldn't be allowed cos it's not his
	// turn!
	err = game.PlayerAction("Steve", "draw", []string{})
	if err == nil {
		t.Fatal(
			"Steve was allowed to draw a card even though it wasn't his turn!")
	}
	// More to come!
}
