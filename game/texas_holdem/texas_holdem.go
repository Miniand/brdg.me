package texas_holdem

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/beefsack/brdg.me/game/card"
	"github.com/beefsack/brdg.me/game/log"
	"github.com/beefsack/brdg.me/game/poker"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	STARTING_MONEY            = 100
	STARTING_MINIMUM_BET      = 10
	HANDS_PER_BLINDS_INCREASE = 5
)

type Game struct {
	Players                  []string
	CurrentPlayer            int
	CurrentDealer            int
	PlayerHands              []card.Deck
	CommunityCards           card.Deck
	Deck                     card.Deck
	Log                      log.Log
	PlayerMoney              []int
	Bets                     []int
	FoldedPlayers            []bool
	MinimumBet               int
	LargestRaise             int
	LastBettingPlayer        int
	HandsSinceBlindsIncrease int
}

func RenderCash(amount int) string {
	return fmt.Sprintf(`{{c "green"}}$%d{{_c}}`, amount)
}

func RenderCashFixedWidth(amount int) string {
	output := RenderCash(amount)
	if amount < 10 {
		output += " "
	}
	if amount < 100 {
		output += " "
	}
	return output
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 || len(players) > 9 {
		return errors.New("Texas hold 'em is limited to 2 - 9 players")
	}
	g.Log = log.NewLog()
	g.Players = players
	g.PlayerHands = make([]card.Deck, len(g.Players))
	g.PlayerMoney = make([]int, len(g.Players))
	for i, _ := range g.Players {
		g.PlayerMoney[i] = STARTING_MONEY
	}
	g.MinimumBet = STARTING_MINIMUM_BET
	// Pick a random starting player
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.CurrentDealer = r.Int() % len(g.Players)
	g.NewHand()
	return nil
}

func (g *Game) NewHand() {
	var (
		smallBlindPlayer, bigBlindPlayer int
	)
	// Reset values
	g.FoldedPlayers = make([]bool, len(g.Players))
	g.Bets = make([]int, len(g.Players))
	g.LargestRaise = 0
	g.NewBettingRound()
	activePlayers := g.ActivePlayers()
	numActivePlayers := len(activePlayers)
	// Raise blinds if we need to
	if g.HandsSinceBlindsIncrease >= HANDS_PER_BLINDS_INCREASE {
		g.HandsSinceBlindsIncrease = 0
		g.MinimumBet *= 2
		g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Minimum bet increased to %s", RenderCash(g.MinimumBet))))
	} else {
		g.HandsSinceBlindsIncrease += 1
	}
	// Set a new active dealer
	g.CurrentDealer = g.NextActivePlayerNumFrom(g.CurrentDealer)
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s is the new dealer",
		g.Players[g.CurrentDealer])))
	// Blinds
	if numActivePlayers == 2 {
		// Special dead-to-head rules for 2 player
		// @see https://en.wikipedia.org/wiki/Texas_hold_'em#Betting_structures
		smallBlindPlayer = g.CurrentDealer
	} else {
		smallBlindPlayer = g.NextActivePlayerNumFrom(g.CurrentDealer)
	}
	g.BetUpTo(smallBlindPlayer, g.MinimumBet/2)
	bigBlindPlayer = g.NextActivePlayerNumFrom(smallBlindPlayer)
	g.BetUpTo(bigBlindPlayer, g.MinimumBet)
	// Make the current player the one next to the big blind
	g.CurrentPlayer = g.NextActivePlayerNumFrom(bigBlindPlayer)
	// Shuffle and deal two cards to each player
	g.CommunityCards = card.Deck{}
	g.Deck = card.Standard52DeckAceHigh().Shuffle()
	for i, _ := range activePlayers {
		g.PlayerHands[i], g.Deck = g.Deck.PopN(2)
		g.PlayerHands[i] = g.PlayerHands[i].Sort()
	}
}

// Remaining players who haven't busted yet
func (g *Game) RemainingPlayers() map[int]string {
	remaining := map[int]string{}
	for i, p := range g.Players {
		if g.PlayerMoney[i] > 0 || g.Bets[i] > 0 {
			remaining[i] = p
		}
	}
	return remaining
}

// Active players are players who are active in the current hand
func (g *Game) ActivePlayers() map[int]string {
	active := map[int]string{}
	for i, p := range g.RemainingPlayers() {
		if !g.FoldedPlayers[i] {
			active[i] = p
		}
	}
	return active
}

func (g *Game) NextActivePlayerNumFrom(playerNum int) int {
	if len(g.ActivePlayers()) == 0 {
		panic("No active players")
	}
	playerNum = (playerNum + 1) % len(g.Players)
	for (g.PlayerMoney[playerNum] == 0 && g.Bets[playerNum] == 0) ||
		g.FoldedPlayers[playerNum] {
		playerNum = (playerNum + 1) % len(g.Players)
	}
	return playerNum
}

func (g *Game) NextRemainingPlayerNumFrom(playerNum int) int {
	if len(g.RemainingPlayers()) == 0 {
		panic("No remaining players")
	}
	playerNum = (playerNum + 1) % len(g.Players)
	for g.PlayerMoney[playerNum] == 0 && g.Bets[playerNum] == 0 {
		playerNum = (playerNum + 1) % len(g.Players)
	}
	return playerNum
}

func (g *Game) BetUpTo(playerNum int, amount int) int {
	betAmount := min(amount, g.PlayerMoney[playerNum])
	err := g.Bet(playerNum, betAmount)
	if err != nil {
		panic(err.Error())
	}
	return betAmount
}

func (g *Game) Bet(playerNum int, amount int) error {
	if g.PlayerMoney[playerNum] < amount {
		return errors.New("Not enough money")
	}
	g.Bets[playerNum] += amount
	g.PlayerMoney[playerNum] -= amount
	g.LargestRaise = max(amount, g.LargestRaise)
	g.LastBettingPlayer = playerNum
	return nil
}

func (g *Game) PlayerAction(player string, action string, args []string) (err error) {
	playerNum, err := g.PlayerNum(player)
	if err == nil {
		switch strings.ToLower(action) {
		case "fold":
			err = g.Fold(playerNum)
		case "call":
			err = g.Call(playerNum)
		case "raise":
			if len(args) == 0 {
				err = errors.New("You must specify an amount to raise")
			} else {
				amount, err := strconv.Atoi(args[0])
				if err != nil {
					err = errors.New("Could not understand your raise amount, only use numbers and no punctuation or symbols")
				} else {
					err = g.Raise(playerNum, amount)
				}
			}
		case "allin":
			err = g.AllIn(playerNum)
		default:
			err = errors.New(fmt.Sprintf("Unknown command: %s", action))
		}
	}
	return
}

func (g *Game) PlayerNum(player string) (int, error) {
	for playerNum, name := range g.Players {
		if player == name {
			return playerNum, nil
		}
	}
	return 0, errors.New("Could not find player with that name")
}

func (g *Game) Fold(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	g.FoldedPlayers[playerNum] = true
	if len(g.ActivePlayers()) == 1 {
		// Everyone folded
		for activePlayerNum, _ := range g.ActivePlayers() {
			g.PlayerMoney[activePlayerNum] += g.Pot()
			g.NewHand()
			return nil
		}
	} else {
		g.NextPlayer()
	}
	return nil
}

func (g *Game) Call(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	difference := g.CurrentBet() - g.Bets[playerNum]
	if g.PlayerMoney[playerNum] < difference {
		return errors.New("You don't have enough to call, you can only go allin")
	}
	err := g.Bet(playerNum, difference)
	if err != nil {
		return err
	}
	g.NextPlayer()
	return nil
}

func (g *Game) Raise(playerNum int, amount int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	minRaise := max(g.MinimumBet, g.LargestRaise)
	difference := g.CurrentBet() - g.Bets[playerNum]
	if amount < minRaise {
		return errors.New(fmt.Sprintf(
			"Your raise must be at least %d", minRaise))
	}
	err := g.Bet(playerNum, difference+amount)
	if err != nil {
		return err
	}
	if amount > g.LargestRaise {

	}
	g.NextPlayer()
	return nil
}

func (g *Game) AllIn(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	err := g.Bet(playerNum, g.PlayerMoney[playerNum])
	if err != nil {
		return err
	}
	g.NextPlayer()
	return nil
}

func (g *Game) NextPlayer() {
	for {
		g.CurrentPlayer = g.NextActivePlayerNumFrom(g.CurrentPlayer)
		if g.CurrentPlayer == g.LastBettingPlayer {
			g.NextPhase()
			break
		}
		if g.PlayerMoney[g.CurrentPlayer] > 0 {
			break
		}
	}
}

func (g *Game) NextPhase() {
	switch len(g.CommunityCards) {
	case 0:
		g.Flop()
	case 3:
		g.Turn()
	case 4:
		g.River()
	case 5:
		g.Showdown()
	}
}

func (g *Game) Flop() {
	g.NewCommunityCards(3)
	g.NewBettingRound()
}

func (g *Game) Turn() {
	g.NewCommunityCards(1)
	g.NewBettingRound()
}

func (g *Game) River() {
	g.NewCommunityCards(1)
	g.NewBettingRound()
}

func (g *Game) Showdown() {
	for g.Pot() > 0 {
		// Find the minimum bet
		smallest := g.SmallestBet()
		pot := 0
		handResults := map[int]poker.HandResult{}
		for playerNum, b := range g.Bets {
			contribution := min(b, smallest)
			pot += contribution
			g.Bets[playerNum] -= contribution
			if !g.FoldedPlayers[playerNum] {
				handResults[playerNum] = poker.Result(g.PlayerHands[playerNum])
			}
		}
		winners := poker.WinningHandResult(handResults)
		potPerPlayer := pot / len(winners)
		for _, winner := range winners {
			g.PlayerMoney[winner] += potPerPlayer
		}
		remainder := pot - potPerPlayer*len(winners)
		if remainder > 0 {
			g.PlayerMoney[g.NextRemainingPlayerNumFrom(g.CurrentDealer)] +=
				remainder
		}
	}
	if !g.IsFinished() {
		g.NewHand()
	}
}

func (g *Game) CurrentBet() int {
	currentBet := 0
	for _, b := range g.Bets {
		if b > currentBet {
			currentBet = b
		}
	}
	return currentBet
}

func (g *Game) Pot() int {
	total := 0
	for _, b := range g.Bets {
		total += b
	}
	return total
}

func (g *Game) SmallestBet() int {
	bet := 0
	firstRun := true
	for playerNum, _ := range g.ActivePlayers() {
		if firstRun || g.Bets[playerNum] < bet {
			bet = g.Bets[playerNum]
			firstRun = false
		}
	}
	return bet
}

func (g *Game) NewCommunityCards(n int) {
	var cards card.Deck
	cards, g.Deck = g.Deck.PopN(3)
	g.CommunityCards = g.CommunityCards.PushMany(cards)
}

func (g *Game) NewBettingRound() {
	g.CurrentPlayer = g.NextActivePlayerNumFrom(g.CurrentDealer)
}

func (g *Game) Name() string {
	return "Texas hold 'em"
}

func (g *Game) Identifier() string {
	return "texas_holdem"
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

func (g *Game) RenderForPlayer(string) (string, error) {
	return "", nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.RemainingPlayers()) < 2
}

func (g *Game) Winners() []string {
	activePlayers := g.ActivePlayers()
	if len(activePlayers) == 1 {
		for _, p := range activePlayers {
			return []string{p}
		}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func min(numbers ...int) int {
	l := len(numbers)
	if l == 0 {
		panic("Requires at least one int")
	}
	m := numbers[0]
	if l > 1 {
		for i := 1; i < l; i++ {
			if numbers[i] < m {
				m = numbers[i]
			}
		}
	}
	return m
}

func max(numbers ...int) int {
	l := len(numbers)
	if l == 0 {
		panic("Requires at least one int")
	}
	m := numbers[0]
	if l > 1 {
		for i := 1; i < l; i++ {
			if numbers[i] > m {
				m = numbers[i]
			}
		}
	}
	return m
}
