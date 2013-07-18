package lost_cities

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"
)

type Game struct {
	// Player 1 is at index 0, player 2 is at index 1
	Players []string
	// Currently moving player, stored as an int to make it easier to reference
	// hands and scores
	CurrentlyMoving int
	// The current phase
	TurnPhase int
	// The board contains the game state of the current round
	Board Board
	// Round out of 3
	Round int
	// Round scores are a multidimensional array, by two players then by three
	// rounds
	RoundScores [2][3]int
	// Tracks which players are ready for the next round before starting the
	// next round, to give players a chance to see the board after a round
	ReadyPlayers [2]bool
}

// The board consists of two players hands, a discard hand, and a draw pile
type Board struct {
	// Player hands are an array of hands, indexed 0 for player 1 and 1 for
	// player 2
	PlayerHands [2]Hand
	// The discard piles is stored as a hand, because it uses the same structure
	DiscardPiles Hand
	// The draw pile is just a flat array of cards because it doesn't need to be
	// grouped into suits
	DrawPile []Card
}

// A hand is cards sorted into suits.  Each player has a hand, and the centre
// discard piles are also stored as a hand.
type Hand struct {
	CardsBySuit [5][]Card
}

// A card is a suit and a value
type Card struct {
	// Suits are defined by the colour consts above 0-4
	Suit int
	// The value is 2-10 for cards an 0 for multiplier
	Value int
}

// Turn phase constant
const (
	// The first half of the turn, playing or discarding a card
	TURN_PHASE_PLAY_OR_DISCARD = iota
	// The second half of the turn, drawing a card from discards or the deck
	TURN_PHASE_DRAW
)

// Suit constants, iota numbers them incrementally from 0, so they are 0-4
const (
	SUIT_RED = iota
	SUIT_GREEN
	SUIT_BLUE
	SUIT_WHITE
	SUIT_YELLOW
)

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Lost Cities requires 2 spieler")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
	err := g.InitRound()
	if err != nil {
		return err
	}
	g.CurrentlyMoving = r.Int() % 2
	return nil
}

// Shuffle cards and deal hands, set the start player, set the turn phase etc
func (g *Game) InitRound() error {
	return nil
}

// Gets the player number, given a string name
func (g *Game) PlayerFromString(player string) (int, error) {
	for key, name := range g.Players {
		if name == player {
			return key, nil
		}
	}
	return 0, errors.New("Couldn't find player with name: " + player)
}

// Takes a string like b6, rx, y10 and turns it into a Card object
func (g *Game) ParseCardString(cardString string) (Card, error) {
	return Card{}, nil
}

func (g *Game) PlayerAction(player, action string, params []string) error {
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	switch strings.ToLower(action) {
	case "play":
		err = g.PlayCard(playerNum, Card{})
	case "discard":
		err = g.DiscardCard(playerNum, Card{})
	case "take":
		err = g.TakeCard(playerNum, SUIT_RED)
	case "draw":
		err = g.DrawCard(playerNum)
	case "ready":
		err = g.PlayerReady(playerNum)
	default:
		err = errors.New("Did not understand your action: " + action)
	}
	return err
}

func (g *Game) Name() string {
	return "Lost Cities"
}

func (g *Game) Identifier() string {
	return "lost_cities"
}

func (g *Game) Encode() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) Decode(data []byte) error {
	return json.Unmarshal(data, g)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{}
}

// Returns the full set of cards in a game, 3 investment cards and 9 point cards
// for each expedition totalling 60 cards
func (g *Game) AllCards() []Card {
	return []Card{}
}

// Play a card from the hand into an expedition, checking that it is the
// player's turn, that they have the card in their hand, and that they are able
// to play the card.  Return an error if any of these don't pass.
func (g *Game) PlayCard(player int, card Card) error {
	return nil
}

// Take a card from discard stacks into the hand, checking that it is the
// player's turn, and that the discard stack has cards in it.  Return an error
// if any of these don't pass.
func (g *Game) TakeCard(player int, suit int) error {
	return nil
}

// Play a card from the hand into an expedition, checking that it is the
// player's turn, that they have the card in their hand, and that they are able
// to play the card.  Return an error if any of these don't pass.
func (g *Game) DrawCard(player int) error {
	return nil
}

// Play a card from the hand into an expedition, checking that it is the
// player's turn, that they have the card in their hand, and that they are able
// to play the card.  Return an error if any of these don't pass.
func (g *Game) DiscardCard(player int, card Card) error {
	return nil
}

// Mark a player as ready for the next round, gives them a chance to see the
// board and scores after the round has ended.  Check that it is actually the
// end of the round and that it isn't the last round already, returning an error
// if either of those fail.
func (g *Game) PlayerReady(player int) error {
	return nil
}
