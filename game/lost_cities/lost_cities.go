package lost_cities

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/beefsack/brdg.me/game/card"
	"math/rand"
	"strings"
	"time"
	"fmt"
	"strconv"
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
	// Player hands are an array of cards, indexed 0 for player 1 and 1 for
	// player 2
	PlayerHands [2]card.Deck
	// Player expeditions are an array of suited piles, indexed 0 for player 1
	// and 1 for player 2
	PlayerExpeditions [2][5]card.Deck
	// The discard piles is stored as a hand, because it uses the same structure
	DiscardPiles [5]card.Deck
	// The draw pile is just a flat array of cards because it doesn't need to be
	// grouped into suits
	DrawPile card.Deck
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

// Suit colours map to ansi colours
// @see http://en.wikipedia.org/wiki/ANSI_escape_code#Colors
var CardColours = map[int]string{
	SUIT_RED:    "red",
	SUIT_GREEN:  "green",
	SUIT_BLUE:   "blue",
	SUIT_WHITE:  "gray",
	SUIT_YELLOW: "yellow",
}

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
	g.Board.DrawPile = g.AllCards().Shuffle()
	g.Board.PlayerHands[0], g.Board.DrawPile = g.Board.DrawPile.PopN(5)
	g.Board.PlayerHands[0] = g.Board.PlayerHands[0].Sort()
	g.Board.PlayerHands[1], g.Board.DrawPile = g.Board.DrawPile.PopN(5)
	g.Board.PlayerHands[1] = g.Board.PlayerHands[1].Sort()
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
func (g *Game) ParseCardString(cardString string) (card.SuitRankCard, error) {
	suitnum:=0
	if len(cardString)<2{
		return card.SuitRankCard{}, errors.New("not lengthy enough (heyoooo!)")
	}
	suit := cardString[0]

	fmt.Println(suit)
	fmt.Println(cardString)	
val,err:=strconv.Atoi(cardString[1:])
if err!=nil{
		return card.SuitRankCard{}, err
	}



switch suit {
			case 'r':
				suitnum = SUIT_RED
			case 'y':
				suitnum = SUIT_YELLOW
			case 'b':
				suitnum = SUIT_BLUE
			case 'w':
				suitnum = SUIT_WHITE
			case 'g':
				suitnum = SUIT_GREEN
			}

	return 	card.SuitRankCard{
			Suit: suitnum,
			Rank: val,
	}, nil


}

func (g *Game) PlayerAction(player, action string, params []string) error {
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	switch strings.ToLower(action) {
	case "play":
		if len(params) == 0 {
			return errors.New("You must specify a card to play, such as r5")
		}
		c, err := g.ParseCardString(params[0])
		if err != nil {
			return err
		}
		err = g.PlayCard(playerNum, c)
	case "discard":
		if len(params) == 0 {
			return errors.New("You must specify a card to play, such as r5")
		}
		c, err := g.ParseCardString(params[0])
		if err != nil {
			return err
		}
		err = g.DiscardCard(playerNum, c)
		g.PlayerAction(player, "draw", []string{})
	case "take":
		// @todo actually detect the suit from params[0], make sure to check
		// @params length > 0
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

func RegisterGobTypes() {
	gob.Register(card.SuitRankCard{})
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(g)
	return buf.Bytes(), err
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(g)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) RenderCard(card card.SuitRankCard) string {
	// @todo Actually do output from card suit and value.  Maybe make sure
	// there's a trailing space if the card value isn't 10, to make sure
	// everything lines up nicely.
	return `{{c "` + CardColours[card.Suit] + `"}}R5{{_c}}`
}

func (g *Game) PlayerList() []string {
	return g.Players
}

// Check that it's the end of the round, that there are no more draw cards left
func (g *Game) IsEndOfRound() bool {
	return false
}

// Check that it's the end of the third round
func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

// Whose turn it is right now, if it's the end of the round this should return
// all players that haven't marked themselves as ready
func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentlyMoving]}
}

// Returns the full set of cards in a game, 3 investment cards and 9 point cards
// for each expedition totalling 60 cards
func (g *Game) AllCards() card.Deck {
	deck := card.Deck{}
	var value int
	for suit := SUIT_RED; suit < SUIT_YELLOW; suit++ {
		for y := 0; y < 12; y++ {
			switch y {
			case 0:
				value = 0
			case 1:
				value = 0
			case 2:
				value = 0
			case 3:
				value = 2
			case 4:
				value = 3
			case 5:
				value = 4
			case 6:
				value = 5
			case 7:
				value = 6
			case 8:
				value = 7
			case 9:
				value = 8
			case 10:
				value = 9
			case 11:
				value = 10
			}
			deck = deck.Push(card.SuitRankCard{
				Suit: suit,
				Rank: value,
			})
		}
	}
	return deck

}

// Play a card from the hand into an expedition, checking that it is the
// player's turn, that they have the card in their hand, and that they are able
// to play the card.  Return an error if any of these don't pass.
func (g *Game) PlayCard(player int, c card.SuitRankCard) error {
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

// Discard a card from the hand into a discard stack, checking that it is the
// player's turn, that they have the card in their hand, 
// Return an error if any of these don't pass.
func (g *Game) DiscardCard(player int, c card.SuitRankCard) error {
	removeCount:=0
	g.Board.PlayerHands[player], removeCount = g.Board.PlayerHands[player].Remove(c,1)
	if removeCount==0{
		return errors.New ("did not have card in hand")
	}
	g.Board.DiscardPiles[c.Suit]=g.Board.DiscardPiles[c.Suit].Push(c)
	return nil
}

// Mark a player as ready for the next round, gives them a chance to see the
// board and scores after the round has ended.  Check that it is actually the
// end of the round and that it isn't the last round already, returning an error
// if either of those fail.  If the player is the last person to be ready, call
// InitRound to start a new round.
func (g *Game) PlayerReady(player int) error {
	return nil
}

// Calculate the current score for this round for a player.
func (g *Game) CurrentRoundPlayerScore(player int) int {
	// @todo You want to be looking at g.Board.PlayerExpeditions[player][SUIT_RED] etc
	return 0
}
