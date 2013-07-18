package lost_cities

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

// Suit constants, iota numbers them incrementally from 0, so they are 0-4
const (
	RED = iota
	GREEN
	BLUE
	WHITE
	YELLOW
)

// A card is a suit and a value
type Card struct {
	// Suits are defined by the colour consts above 0-4
	Suit int
	// The value is 2-10 for cards an 0 for multiplier
	Value int
}

// A hand is cards sorted into suits.  Each player has a hand, and the centre
// discard piles are also stored as a hand.
type Hand struct {
	Suits [5][]Card
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

type Game struct {
	// Player 1 is at index 0, player 2 is at index 1
	Players []string
	// Currently moving player, stored as an int to make it easier to reference
	// hands and scores
	CurrentlyMoving int
	// The board contains the game state of the current round
	Board Board
	// Round out of 3
	Round int
	// Round scores are a multidimensional array, by two players then by three
	// rounds
	RoundScores [2][3]int
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Lost Cities requires 2 spieler")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
	g.CurrentlyMoving = r.Int() % 2
	return nil
}

func (g *Game) PlayerAction(player, action string, params []string) error {
	return nil
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
