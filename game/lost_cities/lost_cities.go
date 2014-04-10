package lost_cities

import (
	"bytes"
	"encoding/gob"
	"errors"
	//"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"math/rand"
	"strconv"
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
	Log         *log.Log
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

var SuitShortNames = map[int]string{
  SUIT_RED: "R",
  SUIT_GREEN:  "G",
  SUIT_BLUE:   "B",
  SUIT_WHITE:  "W",
  SUIT_YELLOW: "Y",
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Lost Cities requires 2 spieler")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
	g.Log = log.New()
	err := g.InitRound()
	if err != nil {
		return err
	}
	g.CurrentlyMoving = r.Int() % 2
	return nil
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

// Shuffle cards and deal hands, set the start player, set the turn phase etc
func (g *Game) InitRound() error {
	g.Board.DrawPile = g.AllCards().Shuffle()
	g.Board.PlayerHands[0], g.Board.DrawPile = g.Board.DrawPile.PopN(8)
	g.Board.PlayerHands[0] = g.Board.PlayerHands[0].Sort()
	g.Board.PlayerHands[1], g.Board.DrawPile = g.Board.DrawPile.PopN(8)
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
	suitnum := 0
	val := 0
	var err error
	// fmt.Println("val")
	// fmt.Println(val)
	if len(cardString) < 2 {
		return card.SuitRankCard{}, errors.New("not lengthy enough (heyoooo!)")
	}
	// fmt.Println("cardstring:")
	// fmt.Println(cardString)
	// fmt.Println("cardstring 1::")
	// fmt.Println(cardString[1:])
	suit := strings.ToLower(cardString[0:1])
	// fmt.Println("suit")
	// fmt.Println(suit)
	if cardString[1:] == "x" {
		val = 0
	} else {
		val, err = strconv.Atoi(cardString[1:])
		if err != nil {
			return card.SuitRankCard{}, err
		}
		// fmt.Println("val")
		// fmt.Println(val)
	}
	// fmt.Println("val now")
	// fmt.Println(val)
	switch suit {
	case "r":
		suitnum = SUIT_RED
	case "y":
		suitnum = SUIT_YELLOW
	case "b":
		suitnum = SUIT_BLUE
	case "w":
		suitnum = SUIT_WHITE
	case "g":
		suitnum = SUIT_GREEN
	default:
		return card.SuitRankCard{}, errors.New("Could not parse suit")
	}

	return card.SuitRankCard{
		Suit: suitnum,
		Rank: val,
	}, nil
}

// Defines which commands are available for Lost Cities, see the _command.go
// files in this directory.
func (g *Game) Commands() []command.Command {
	return []command.Command{
		DiscardCommand{},
		PlayCommand{},
		DrawCommand{},
		TakeCommand{},
	}
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
	output := bytes.NewBufferString("")
	output.WriteString("player:\n")
	output.WriteString(player + "\n")
	val, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	output.WriteString("your hand:\n")
	for count := 0; count < len(g.Board.PlayerHands[val]); count++ {

		output.WriteString(g.RenderCard(g.Board.PlayerHands[val][count].(card.SuitRankCard)) + "\n")

	}
	output.WriteString("your expeditions:\n")
	for suits := 0; suits < 5; suits++ {
		for count := 0; count < len(g.Board.PlayerExpeditions[val][suits]); count++ {
			output.WriteString(g.RenderCard(g.Board.PlayerExpeditions[val][suits][count].(card.SuitRankCard)) + "\n")

		}
	}
	opp := 0
	if val == 0 {
		opp = 1
	}
	output.WriteString("the other dude's expeditions:\n")
	for suits := 0; suits < 5; suits++ {
		for count := 0; count < len(g.Board.PlayerExpeditions[opp][suits]); count++ {
			output.WriteString(g.RenderCard(g.Board.PlayerExpeditions[opp][suits][count].(card.SuitRankCard)) + "\n")

		}
	}
	output.WriteString("the discard piles:\n")
	for suits := 0; suits < 5; suits++ {
		if len(g.Board.DiscardPiles[suits])>0 {
			output.WriteString(g.RenderCard(g.Board.DiscardPiles[suits][len(g.Board.DiscardPiles[suits])-1].(card.SuitRankCard)) + "\n")

		}
	}
	return output.String(), nil

}

func (g *Game) RenderCard(card card.SuitRankCard) string {
	// @todo Actually do output from card suit and value.  Maybe make sure
	// there's a trailing space if the card value isn't 10, to make sure
	// everything lines up nicely.
	return `{{c "` + CardColours[card.Suit] + `"}}` + strconv.Itoa(card.Rank) + `{{_c}}` + `: ` + SuitShortNames[card.Suit] + strconv.Itoa(card.Rank)
	return CardColours[card.Suit] + " " + strconv.Itoa(card.Rank)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

// Check that it's the end of the third round
func (g *Game) IsFinished() bool {
	if g.Round == 3 {
		return true
	}
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

// Whose turn it is right now
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
	removeCount := 0
	//fmt.Println("in PlayCard c.Suit")
	//fmt.Println(c.Suit)
	//fmt.Println("in PlayCard c.Rank")
	//fmt.Println(c.Rank)
	//fmt.Println("%#v\n", g.Board.PlayerExpeditions[1])
	//fmt.Println(" \n")

	//len(g.Board.PlayerExpeditions[player][c.Suit])-1
	if len(g.Board.PlayerExpeditions[player][c.Suit]) > 0 {
		if c.Rank < g.Board.PlayerExpeditions[player][c.Suit][len(g.Board.PlayerExpeditions[player][c.Suit])-1].(card.SuitRankCard).Rank {
			return errors.New("too low!")
		}
	}

	g.Board.PlayerHands[player], removeCount = g.Board.PlayerHands[player].Remove(c, 1)

	//fmt.Println("%#v\n", g.Board.PlayerExpeditions[1])
	//fmt.Println(" \n")
	//fmt.Println(g.Board.PlayerExpeditions[player][c.Suit][0])
	if removeCount == 0 {
		return errors.New("did not have card in hand")
	}
	g.Board.PlayerExpeditions[player][c.Suit] = g.Board.PlayerExpeditions[player][c.Suit].Push(c)
	g.TurnPhase = 1
	return nil
}

// Take a card from discard stacks into the hand, checking that it is the
// player's turn, and that the discard stack has cards in it.  Return an error
// if any of these don't pass.
func (g *Game) TakeCard(player int, suit int) error {
	// fmt.Println("gonna take a card")
	// fmt.Println(suit)
	if len(g.Board.DiscardPiles[suit]) == 0 {
		return errors.New("There are no cards in that discard pile")
	}

	var drawnCard card.Card
	drawnCard, g.Board.DiscardPiles[suit] = g.Board.DiscardPiles[suit].Pop()
	// fmt.Println(drawnCard)

	g.Board.PlayerHands[player] = g.Board.PlayerHands[player].Push(drawnCard)
	// fmt.Println(g.Board.PlayerHands[player])
	//fmt.Println(c.Suit)
	//fmt.Println("in PlayCard c.Rank")
	//fmt.Println(c.Rank)

	//g.Board.DiscardPiles[suit], removeCount = g.Board.PlayerHands[player].Remove(c, 1)

	//fmt.Println("%#v\n", g.Board.PlayerExpeditions[1])
	//fmt.Println(" \n")
	//fmt.Println(g.Board.PlayerExpeditions[player][c.Suit][0])
	//if removeCount==0{
	//	return errors.New ("did not have card in hand")
	//}
	if g.CurrentlyMoving == 1 {
		g.CurrentlyMoving = 0
	} else {
		g.CurrentlyMoving = 1
	}
	g.TurnPhase = 0

	return nil

}

// Take a card from the draw pile into the hand, checking that it is the
// player's turn and that there are cards in the stack.
// Return an error if any of these don't pass.
func (g *Game) DrawCard(player int) error {

	var drawnCard card.Card
	drawnCard, g.Board.DrawPile = g.Board.DrawPile.Pop()
	// Then put it into the player's hand
	g.Board.PlayerHands[player] = g.Board.PlayerHands[player].Push(drawnCard)
	if g.CurrentlyMoving == 1 {
		g.CurrentlyMoving = 0
	} else {
		g.CurrentlyMoving = 1
	}
	g.TurnPhase = 0
	if len(g.Board.DrawPile) == 0 {
		g.Round = g.Round + 1

	}
	return nil
}

// Discard a card from the hand into a discard stack, checking that it is the
// player's turn, that they have the card in their hand,
// Return an error if any of these don't pass.
func (g *Game) DiscardCard(player int, c card.SuitRankCard) error {
	removeCount := 0
	g.Board.PlayerHands[player], removeCount = g.Board.PlayerHands[player].Remove(c, 1)
	if removeCount == 0 {
		return errors.New("did not have card in hand")
	}
	g.Board.DiscardPiles[c.Suit] = g.Board.DiscardPiles[c.Suit].Push(c)
	g.TurnPhase = 1

	//fmt.Println("I get to here")
	return nil
}

// Calculate the current score for this round for a player.
func (g *Game) CurrentRoundPlayerScore(player int) int {
	// @todo You want to be looking at g.Board.PlayerExpeditions[player][SUIT_RED] etc
	return 0
}

// Run through an expedition to calculate the score.  Ignore suits here, just
// focus on values, the rest of the game logic can ensure the deck is of the
// right suit.
func ScoreExpedition(hand card.Deck) int {
	//if hand!=null{
	// fmt.Println(hand)
	// fmt.Println("array length:")
	// fmt.Println(len(hand))
	total := 0
	if len(hand) != 0 {
		total = -20
	}

	investments := 0
	//times by number of investments+1
	for count := 0; count < len(hand); count++ {
		// fmt.Println(hand[count].(card.SuitRankCard).Rank)
		if (hand[count].(card.SuitRankCard).Rank) == 0 {
			investments++
		} else {
			total = total + hand[count].(card.SuitRankCard).Rank
		}

	}
	total = total * (investments + 1)
	if len(hand) >= 8 {
		total = total + 20
	}
	return total
}
