package lost_cities

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
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
	// Mick is the first player
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
	// Just set the draw pile to have a couple of cards so we can finish the
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
	//t.Logf("%#v\n", game.Board.PlayerExpeditions[1])
	//t.Logf(" \n")
	if game.IsEndOfRound() || game.IsFinished() {
		t.Fatal("Why is it the end of the round if we've just started?")
	}

	// MICK FIRST TURN
	// Play or discard

	if game.TurnPhase != TURN_PHASE_PLAY_OR_DISCARD {
		t.Fatal("The turn phase isn't for the player to play or discard")
	}
	// Mick discards red 5
	_, err := command.CallInCommands("Mick", game, "discard r5", game.Commands())
	if err != nil {
		t.Fatal(err)
	}
	// Let's check to make sure it actually happened
	if len(game.Board.PlayerHands[0]) != 7 {
		t.Fatal("Red 5 wasn't removed from Mick's hand when he tried to discard it")
	}
	if len(game.Board.DiscardPiles[SUIT_RED]) != 1 {
		t.Fatal("The red discard pile doesn't have any cards in it")
	}
	firstRed := game.Board.DiscardPiles[SUIT_RED][0].(card.SuitRankCard)
	if firstRed.Suit != SUIT_RED && firstRed.Rank != 5 {
		t.Fatal("Red 5 wasn't discard onto the red discard pile")
	}

	// MICK FIRST TURN
	// Draw

	if game.TurnPhase != TURN_PHASE_DRAW {
		t.Fatal("The turn phase didn't change to DRAW")
	}
	// Steve tries to butt in, but he shouldn't be allowed cos it's not his
	// turn!
	_, err = command.CallInCommands("Steve", game, "draw", game.Commands())
	if err == nil {
		t.Fatal(
			"Steve was allowed to draw a card even though it wasn't his turn!")
	}
	// Mick draws from the draw pile
	_, err = command.CallInCommands("Mick", game, "draw", game.Commands())
	if err != nil {
		t.Fatal(err)
	}
	if len(game.Board.PlayerHands[0]) != 8 {
		t.Fatal("Mick's hand isn't 8 cards after drawing")
	}
	if len(game.Board.DrawPile) != 2 {
		t.Fatal("The draw pile didn't reduce by one after drawing")
	}

	// STEVE FIRST TURN
	// Play or discard

	// Make sure the turn changed to Steve
	if game.CurrentlyMoving != 1 {
		t.Fatal("Turn didn't change to Steve since Mick finished playing")
	}
	if game.TurnPhase != TURN_PHASE_PLAY_OR_DISCARD {
		t.Fatal("The turn phase isn't to play or discard")
	}
	// Try to draw first and make sure we aren't allowed
	_, err = command.CallInCommands("Steve", game, "draw", game.Commands())
	if err == nil {
		t.Fatal("The game let Steve draw, he hasn't played yet!")
	}
	//t.Logf("%#v\n", game.Board.PlayerExpeditions[1])
	//t.Logf(" \n")
	// Play a blue 9 and check it actually happened
	_, err = command.CallInCommands("Steve", game, "play B9", game.Commands())
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("%#v\n", game.Board.PlayerExpeditions[1])
	// t.Logf(" \n")
	if len(game.Board.PlayerExpeditions[1][SUIT_BLUE]) != 1 ||
		game.Board.PlayerExpeditions[1][SUIT_BLUE][0].(card.SuitRankCard).Rank != 9 {
		t.Fatal("We couldn't find the blue 9 in Steve's blue player expedition")
	}
	if len(game.Board.PlayerHands[1]) != 7 {
		t.Fatal("Steve's hand wasn't reduced to 7")
	}

	// STEVE FIRST TURN
	// Draw

	// Steve will draw from the red discard pile instead of the draw pile
	_, err = command.CallInCommands("Steve", game, "take r", game.Commands())
	if err != nil {
		t.Fatal(err)
	}
	// Make sure Steve actually took it
	if len(game.Board.DiscardPiles[SUIT_RED]) != 0 {
		t.Fatal("The red discard pile isn't empty after Steve drew from it")
	}
	if len(game.Board.PlayerHands[1]) != 8 {
		t.Fatal("Steve's hand isn't 8 after taking a red card")
	}
	takenCard := game.Board.PlayerHands[1][7].(card.SuitRankCard)
	if takenCard.Suit != SUIT_RED || takenCard.Rank != 5 {
		t.Fatal("The card Steve took into his hand wasn't red 5")
	}

	// MICK SECOND TURN
	// Play or discard

	// Mick will play the yellow investment card he has
	_, err = command.CallInCommands("Mick", game, "play yx", game.Commands())
	if err != nil {
		t.Fatal(err)
	}

	// More to come!
}

func TestExpeditionScores(t *testing.T) {
	var (
		expedition    card.Deck
		expectedScore int
	)

	// Empty expedition should be 0, -20 is not applied
	expedition = card.Deck{}
	expectedScore = 0
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// But if we've got a card, -20 is applied.  This test also tests multiplier
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
	}
	expectedScore = -40
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// Try a hand without a multiplier but only value cards
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 3,
		},
		card.SuitRankCard{
			Rank: 5,
		},
		card.SuitRankCard{
			Rank: 7,
		},
	}
	expectedScore = -5
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// Try a hand with both
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 5,
		},
		card.SuitRankCard{
			Rank: 7,
		},
		card.SuitRankCard{
			Rank: 10,
		},
	}
	expectedScore = 4
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// Check that the 20 point bonus is applied for 8 cards after multiplier is
	// applied
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 2,
		},
		card.SuitRankCard{
			Rank: 3,
		},
		card.SuitRankCard{
			Rank: 4,
		},
		card.SuitRankCard{
			Rank: 5,
		},
		card.SuitRankCard{
			Rank: 6,
		},
		card.SuitRankCard{
			Rank: 7,
		},
		card.SuitRankCard{
			Rank: 10,
		},
	}
	expectedScore = 54
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// But make sure there is no 20 point bonus for only 7 cards!
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 2,
		},
		card.SuitRankCard{
			Rank: 3,
		},
		card.SuitRankCard{
			Rank: 4,
		},
		card.SuitRankCard{
			Rank: 5,
		},
		card.SuitRankCard{
			Rank: 6,
		},
		card.SuitRankCard{
			Rank: 7,
		},
	}
	expectedScore = 14
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// Min score
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 0,
		},
	}
	expectedScore = -80
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)

	// Max score, including 20 point bonus after multiplier
	expedition = card.Deck{
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 0,
		},
		card.SuitRankCard{
			Rank: 2,
		},
		card.SuitRankCard{
			Rank: 3,
		},
		card.SuitRankCard{
			Rank: 4,
		},
		card.SuitRankCard{
			Rank: 5,
		},
		card.SuitRankCard{
			Rank: 6,
		},
		card.SuitRankCard{
			Rank: 7,
		},
		card.SuitRankCard{
			Rank: 8,
		},
		card.SuitRankCard{
			Rank: 9,
		},
		card.SuitRankCard{
			Rank: 10,
		},
	}
	expectedScore = 156
	ensureExpeditionMatchesExpectedScore(t, expedition, expectedScore)
}
func ensureExpeditionMatchesExpectedScore(t *testing.T, expedition card.Deck,
	expectedScore int) {
	actualScore := ScoreExpedition(expedition)
	if ScoreExpedition(expedition) != expectedScore {
		ranks := []int{} // We build raw ranks to make output easier to read
		for _, c := range expedition {
			ranks = append(ranks, c.(card.SuitRankCard).Rank)
		}
		t.Fatalf("Score for expedition %v was expected to be %d, but got %d",
			ranks, expectedScore, actualScore)
	}
}
