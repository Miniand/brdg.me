package lost_cities

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
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
	// LastDiscardedSuit has the suit of the last discarded card, to prevent
	// a player taking the card they just discarded.
	LastDiscardedSuit int
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

const (
	DIR_ASC = iota
	DIR_DESC
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
	SUIT_RED:    "R",
	SUIT_GREEN:  "G",
	SUIT_BLUE:   "B",
	SUIT_WHITE:  "W",
	SUIT_YELLOW: "Y",
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Lost Cities requires 2 spieler")
	}
	g.Players = players
	g.Log = log.New()
	err := g.InitRound()
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

// Shuffle cards and deal hands, set the start player, set the turn phase etc
func (g *Game) InitRound() error {
	g.Board = Board{}
	g.Board.DrawPile = g.AllCards().Shuffle()
	g.Board.PlayerHands[0], g.Board.DrawPile = g.Board.DrawPile.PopN(8)
	g.Board.PlayerHands[0] = g.Board.PlayerHands[0].Sort()
	g.Board.PlayerHands[1], g.Board.DrawPile = g.Board.DrawPile.PopN(8)
	g.Board.PlayerHands[1] = g.Board.PlayerHands[1].Sort()
	p0s := g.PreviousRoundsPlayerScore(0)
	p1s := g.PreviousRoundsPlayerScore(1)
	if p0s > p1s {
		g.CurrentlyMoving = 0
	} else if p1s > p0s {
		g.CurrentlyMoving = 1
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		g.CurrentlyMoving = r.Int() % 2
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(`{{b}}%s{{_b}} will start round %d`,
		render.PlayerName(g.CurrentlyMoving, g.Players[g.CurrentlyMoving]),
		g.Round+1)))
	return nil
}

// Gets the player number, given a string name
func (g *Game) PlayerFromString(player string) (int, error) {
	for key, name := range g.Players {
		if name == player {
			return key, nil
		}
	}
	return 0, fmt.Errorf("couldn't find player with name %s", player)
}

// ParseCardString takes a string like b6, rx, y10 and turns it into a Card
// object.
func (g *Game) ParseCardString(cardString string) (c card.SuitRankCard, err error) {
	if len(cardString) < 2 {
		err = errors.New("card string should be at least 2 characters long")
		return
	}
	if strings.ToUpper(cardString[1:]) == "X" {
		c.Rank = 0
	} else {
		c.Rank, err = strconv.Atoi(cardString[1:])
		if err != nil {
			return
		}
	}
	suit := strings.ToUpper(cardString[:1])
	c.Suit = -1
	for s := SUIT_RED; s <= SUIT_YELLOW; s++ {
		if suit == SuitShortNames[s] {
			c.Suit = s
			break
		}
	}
	if c.Suit == -1 {
		err = errors.New("could not parse suit")
	}
	return
}

// Defines which commands are available for Lost Cities, see the _command.go
// files in this directory.
func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
		DiscardCommand{},
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

func (g *Game) PlayerExpeditionCells(pNum, dir int, title string) (cells [][]string) {
	maxExpSize := 0
	for s := SUIT_RED; s <= SUIT_YELLOW; s++ {
		l := len(g.Board.PlayerExpeditions[pNum][s])
		if l > maxExpSize {
			maxExpSize = l
		}
	}
	for i := 0; i < maxExpSize; i++ {
		row := make([]string, SUIT_YELLOW+2)
		header := ""
		if i == 0 {
			header = fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, title)
		}
		row[0] = header
		for s := SUIT_RED; s <= SUIT_YELLOW; s++ {
			cell := ""
			if len(g.Board.PlayerExpeditions[pNum][s]) > i {
				cell = g.RenderCard(
					g.Board.PlayerExpeditions[pNum][s][i].(card.SuitRankCard))
			}
			row[s+1] = cell
		}
		if dir == DIR_ASC {
			cells = append(cells, row)
		} else {
			cells = append([][]string{row}, cells...)
		}
	}
	return
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	output := bytes.NewBufferString("")
	pNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	// Board
	cells := [][]string{}
	// Opponent area
	cells = append(cells, g.PlayerExpeditionCells((pNum+1)%2, DIR_DESC, "   Them ")...)
	// Discard area
	discard := make([]string, SUIT_YELLOW+2)
	discard[0] = `{{c "gray"}}Discard {{_c}}`
	for s := SUIT_RED; s <= SUIT_YELLOW; s++ {
		var cell string
		l := len(g.Board.DiscardPiles[s])
		if l > 0 {
			cell = g.RenderCard(g.Board.DiscardPiles[s][l-1].(card.SuitRankCard))
		} else {
			cell = fmt.Sprintf(`{{c "%s"}}--{{_c}}`, CardColours[s])
		}
		discard[s+1] = cell
	}
	cells = append(cells, []string{}, discard, []string{})
	// Your area
	cells = append(cells, g.PlayerExpeditionCells(pNum, DIR_ASC, "    You ")...)
	table, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString(table)
	output.WriteString("\n\n")
	// Remaining draw cards
	output.WriteString(fmt.Sprintf("{{b}}Draw cards:{{_b}} %d\n\n",
		len(g.Board.DrawPile)))
	// Your hand
	output.WriteString("{{b}}Your hand:{{_b}}\n")
	hand := []string{}
	for _, c := range g.Board.PlayerHands[pNum] {
		hand = append(hand, g.RenderCard(c.(card.SuitRankCard)))
	}
	output.WriteString(strings.Join(hand, " "))
	output.WriteString("\n\n")
	// Round scores
	cells = [][]string{
		[]string{
			"{{b}}Player{{_b}}",
			"{{b}}R1{{_b}}",
			"{{b}}R2{{_b}}",
			"{{b}}R3{{_b}}",
			"{{b}}Tot{{_b}}",
		},
	}
	for p := 0; p <= 1; p++ {
		row := []string{
			render.PlayerName(p, g.Players[p]),
		}
		for r := 0; r < 3; r++ {
			score := ""
			if g.Round > r {
				score = strconv.Itoa(g.RoundScores[p][r])
			}
			row = append(row, score)
		}
		row = append(row, strconv.Itoa(g.PreviousRoundsPlayerScore(p)))
		cells = append(cells, row)
	}
	table, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString(table)
	return output.String(), nil
}

func (g *Game) RenderCard(card card.SuitRankCard) string {
	rank := strconv.Itoa(card.Rank)
	if rank == "0" {
		rank = "X"
	}
	return fmt.Sprintf(`{{c "%s"}}%s%s{{_c}}`, CardColours[card.Suit],
		SuitShortNames[card.Suit], rank)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

// Check that it's the end of the third round
func (g *Game) IsFinished() bool {
	return g.Round >= 3
}

func (g *Game) Winners() (winners []string) {
	if !g.IsFinished() {
		return
	}
	p0s := g.PreviousRoundsPlayerScore(0)
	p1s := g.PreviousRoundsPlayerScore(1)
	if p0s >= p1s {
		winners = append(winners, g.Players[0])
	}
	if p1s >= p0s {
		winners = append(winners, g.Players[1])
	}
	return
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
	l := len(g.Board.PlayerExpeditions[player][c.Suit])
	if l > 0 && c.Rank <
		g.Board.PlayerExpeditions[player][c.Suit][l-1].(card.SuitRankCard).Rank {
		return errors.New("that card is too low to fit in the expedition")
	}
	g.Board.PlayerHands[player], removeCount = g.Board.PlayerHands[player].Remove(c, 1)
	if removeCount == 0 {
		return errors.New("you do not have that card in hand")
	}
	g.Board.PlayerExpeditions[player][c.Suit] =
		g.Board.PlayerExpeditions[player][c.Suit].Push(c)
	g.TurnPhase = 1
	g.LastDiscardedSuit = -1
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} played {{b}}%s{{_b}}`,
		render.PlayerName(player, g.Players[player]), g.RenderCard(c))))
	return nil
}

// Take a card from discard stacks into the hand, checking that it is the
// player's turn, and that the discard stack has cards in it.  Return an error
// if any of these don't pass.
func (g *Game) TakeCard(player int, suit int) error {
	if len(g.Board.DiscardPiles[suit]) == 0 {
		return errors.New("there are no cards in that discard pile")
	}
	if suit == g.LastDiscardedSuit {
		return errors.New("you can't take the card that you just discarded")
	}
	var drawnCard card.Card
	drawnCard, g.Board.DiscardPiles[suit] = g.Board.DiscardPiles[suit].Pop()
	g.Board.PlayerHands[player] = g.Board.PlayerHands[player].Push(drawnCard).Sort()
	g.NextPlayer()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} took {{b}}%s{{_b}}`,
		render.PlayerName(player, g.Players[player]),
		g.RenderCard(drawnCard.(card.SuitRankCard)))))
	return nil
}

func (g *Game) NextPlayer() {
	g.CurrentlyMoving = (g.CurrentlyMoving + 1) % 2
	g.TurnPhase = TURN_PHASE_PLAY_OR_DISCARD
}

// Take a card from the draw pile into the hand, checking that it is the
// player's turn and that there are cards in the stack.
// Return an error if any of these don't pass.
func (g *Game) DrawCard(player int) error {
	var drawnCard card.Card
	drawnCard, g.Board.DrawPile = g.Board.DrawPile.Pop()
	// Then put it into the player's hand
	g.Board.PlayerHands[player] = g.Board.PlayerHands[player].Push(drawnCard).Sort()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} drew a card, {{b}}%d{{_b}} remaining`,
		render.PlayerName(player, g.Players[player]), len(g.Board.DrawPile))))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		`You drew {{b}}%s{{_b}}`, g.RenderCard(drawnCard.(card.SuitRankCard))),
		[]string{g.Players[g.CurrentlyMoving]}))
	g.NextPlayer()
	if len(g.Board.DrawPile) == 0 {
		g.EndRound()
	}
	return nil
}

func (g *Game) EndRound() {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"{{b}}It is the end of round %d{{_b}}\n", g.Round+1)))
	for p, _ := range g.Players {
		g.RoundScores[p][g.Round] = g.CurrentRoundPlayerScore(p)
	}
	g.Round = g.Round + 1
	if g.Round < 3 {
		g.InitRound()
	}
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
	g.LastDiscardedSuit = c.Suit
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} discarded {{b}}%s{{_b}}`,
		render.PlayerName(player, g.Players[player]), g.RenderCard(c))))
	return nil
}

func (g *Game) PreviousRoundsPlayerScore(player int) int {
	sum := 0
	for _, s := range g.RoundScores[player] {
		sum += s
	}
	return sum
}

// Calculate the current score for this round for a player.
func (g *Game) CurrentRoundPlayerScore(player int) int {
	score := 0
	for s := SUIT_RED; s <= SUIT_YELLOW; s++ {
		score += ScoreExpedition(g.Board.PlayerExpeditions[player][s])
	}
	return score
}

// Run through an expedition to calculate the score.  Ignore suits here, just
// focus on values, the rest of the game logic can ensure the deck is of the
// right suit.
func ScoreExpedition(hand card.Deck) int {
	total := 0
	if len(hand) != 0 {
		total = -20
	}

	investments := 0
	//times by number of investments+1
	for count := 0; count < len(hand); count++ {
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
