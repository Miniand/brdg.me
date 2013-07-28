package texas_holdem

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/beefsack/brdg.me/game/card"
	"github.com/beefsack/brdg.me/game/log"
	"github.com/beefsack/brdg.me/game/poker"
	"github.com/beefsack/brdg.me/render"
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
	LastRaisingPlayer        int
	HandsSinceBlindsIncrease int
}

func RenderCash(amount int) string {
	return fmt.Sprintf(`{{b}}{{c "green"}}$%d{{_c}}{{_b}}`, amount)
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
		g.RenderPlayerName(g.CurrentDealer))))
	// Blinds
	if numActivePlayers == 2 {
		// Special head-to-head rules for 2 player
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
	g.LastRaisingPlayer = g.CurrentPlayer
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
	raiseAmount := amount - g.CurrentBet()
	g.Bets[playerNum] += amount
	g.PlayerMoney[playerNum] -= amount
	g.LargestRaise = max(raiseAmount, g.LargestRaise)
	if raiseAmount > 0 {
		g.LastRaisingPlayer = playerNum
	}
	return nil
}

func (g *Game) PlayerAction(player string, action string, args []string) (err error) {
	playerNum, err := g.PlayerNum(player)
	if err == nil {
		switch strings.ToLower(action) {
		case "check":
			err = g.Check(playerNum)
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

func (g *Game) Check(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	if g.CurrentBet() > g.Bets[playerNum] {
		return errors.New("Cannot check because you are below the bet")
	}
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s checked",
		g.RenderPlayerName(playerNum))))
	g.NextPlayer()
	return nil
}

func (g *Game) Fold(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s folded",
		g.RenderPlayerName(playerNum))))
	g.FoldedPlayers[playerNum] = true
	if len(g.ActivePlayers()) == 1 {
		// Everyone folded
		for activePlayerNum, _ := range g.ActivePlayers() {
			g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"%s won %s", g.RenderPlayerName(activePlayerNum),
				RenderCash(g.Pot()-g.Bets[activePlayerNum]))))
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
	if difference <= 0 {
		return errors.New(
			"You are already at the current bet, you may check if you don't want to raise")
	}
	err := g.Bet(playerNum, difference)
	if err != nil {
		return err
	}
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s called",
		g.RenderPlayerName(playerNum))))
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
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s raised by %s",
		g.RenderPlayerName(playerNum), RenderCash(amount))))
	g.NextPlayer()
	return nil
}

func (g *Game) AllIn(playerNum int) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum {
		return errors.New("Not your turn")
	}
	amount := g.PlayerMoney[playerNum]
	err := g.Bet(playerNum, amount)
	if err != nil {
		return err
	}
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s went all in with %s",
		g.RenderPlayerName(playerNum), RenderCash(amount))))
	g.NextPlayer()
	return nil
}

func (g *Game) NextPlayer() {
	for {
		g.CurrentPlayer = g.NextActivePlayerNumFrom(g.CurrentPlayer)
		if g.CurrentPlayer == g.LastRaisingPlayer {
			g.NextPhase()
			break
		}
		if g.PlayerMoney[g.CurrentPlayer] > 0 {
			break
		}
	}
}

func (g *Game) NextPhase() {
	playerCountWithMoney := 0
	for i, _ := range g.ActivePlayers() {
		if g.PlayerMoney[i] > 0 {
			playerCountWithMoney++
		}
	}
	switch len(g.CommunityCards) {
	case 0:
		g.Flop()
		if playerCountWithMoney < 2 {
			g.NextPhase()
		}
	case 3:
		g.Turn()
		if playerCountWithMoney < 2 {
			g.NextPhase()
		}
	case 4:
		g.River()
		if playerCountWithMoney < 2 {
			g.NextPhase()
		}
	case 5:
		g.Showdown()
	}
}

func (g *Game) Flop() {
	g.NewCommunityCards(3)
	g.Log = g.Log.Add(log.NewPublicMessage("Flop cards are {{b}}" +
		strings.Join(RenderCards(g.CommunityCards), " ") + "{{_b}}"))
	g.NewBettingRound()
}

func (g *Game) Turn() {
	g.NewCommunityCards(1)
	g.Log = g.Log.Add(log.NewPublicMessage("Turn card is {{b}}" +
		g.CommunityCards[3].(card.SuitRankCard).RenderStandard52() + "{{_b}}"))
	g.NewBettingRound()
}

func (g *Game) River() {
	g.NewCommunityCards(1)
	g.Log = g.Log.Add(log.NewPublicMessage("River card is {{b}}" +
		g.CommunityCards[4].(card.SuitRankCard).RenderStandard52() + "{{_b}}"))
	g.NewBettingRound()
}

func (g *Game) Showdown() {
	buf := bytes.NewBufferString("{{b}}Showdown{{_b}}\n")
	for g.Pot() > 0 {
		// Find the minimum bet
		smallest := g.SmallestBet()
		pot := 0
		handResults := map[int]poker.HandResult{}
		handsTable := [][]string{}
		for playerNum, b := range g.Bets {
			if b == 0 {
				continue
			}
			contribution := min(b, smallest)
			pot += contribution
			g.Bets[playerNum] -= contribution
			if !g.FoldedPlayers[playerNum] {
				handResults[playerNum] = poker.Result(
					g.PlayerHands[playerNum].PushMany(g.CommunityCards))
				handsTableRow := []string{g.RenderPlayerName(playerNum)}
				fmt.Println(handResults[playerNum].Cards)
				handsTableRow = append(handsTableRow, strings.Join(
					RenderCards(handResults[playerNum].Cards), " "))
				handsTableRow = append(handsTableRow,
					handResults[playerNum].Name)
				handsTable = append(handsTable, handsTableRow)
			}
		}
		handsTableOutput, err := render.Table(handsTable, 0, 1)
		if err != nil {
			panic(err.Error())
		}
		buf.WriteString(fmt.Sprintf("Pot %s\n%s\n", RenderCash(pot),
			handsTableOutput))
		winners := poker.WinningHandResult(handResults)
		potPerPlayer := pot / len(winners)
		for _, winner := range winners {
			buf.WriteString(fmt.Sprintf("%s won %s (%s)\n",
				g.RenderPlayerName(winner), RenderCash(potPerPlayer),
				handResults[winner].Name))
			g.PlayerMoney[winner] += potPerPlayer
		}
		remainder := pot - potPerPlayer*len(winners)
		if remainder > 0 {
			remainderPlayer := g.NextRemainingPlayerNumFrom(g.CurrentDealer)
			buf.WriteString(fmt.Sprintf("%s took %s due to uneven split",
				g.RenderPlayerName(remainderPlayer), RenderCash(remainder)))
			g.PlayerMoney[remainderPlayer] += remainder
		}
	}
	g.Log = g.Log.Add(log.NewPublicMessage(buf.String()))
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
	cards, g.Deck = g.Deck.PopN(n)
	g.CommunityCards = g.CommunityCards.PushMany(cards)
}

func (g *Game) NewBettingRound() {
	g.CurrentPlayer = g.NextActivePlayerNumFrom(g.CurrentDealer)
	g.LastRaisingPlayer = g.CurrentPlayer
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

func (g *Game) RenderPlayerName(playerNum int) string {
	return render.PlayerName(playerNum, g.Players[playerNum])
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	// Log
	newMessages := g.Log.NewMessagesFor(player)
	if len(newMessages) > 0 {
		buf.WriteString("{{b}}Since last time:{{_b}}\n")
		buf.WriteString(log.RenderMessages(newMessages))
		buf.WriteString("\n\n")
	}
	// Table
	buf.WriteString("{{b}}Community cards{{_b}}:  ")
	buf.WriteString(strings.Join(RenderCards(g.CommunityCards), " "))
	buf.WriteString("\n")
	buf.WriteString("{{b}}Current pot{{_b}}:      ")
	buf.WriteString(RenderCash(g.Pot()))
	buf.WriteString("\n\n")
	// Player specific
	buf.WriteString("{{b}}Your cards{{_b}}:  ")
	buf.WriteString(strings.Join(RenderCards(g.PlayerHands[playerNum]), " "))
	buf.WriteString("\n")
	buf.WriteString("{{b}}Your cash{{_b}}:   ")
	buf.WriteString(RenderCash(g.PlayerMoney[playerNum]))
	buf.WriteString("\n\n")
	if playerNum == g.CurrentPlayer {
		// Available actions
		actions := []string{}
		currentBet := g.CurrentBet()
		currentBetDiff := currentBet - g.Bets[playerNum]
		if currentBetDiff == 0 {
			actions = append(actions, "{{b}}check{{_b}}")
		}
		if currentBetDiff > 0 && currentBetDiff < g.PlayerMoney[playerNum] {
			actions = append(actions, fmt.Sprintf("{{b}}call{{_b}} (for $%d)",
				currentBet-g.Bets[playerNum]), "{{b}}fold{{_b}}")
		}
		if currentBetDiff+g.LargestRaise < g.PlayerMoney[playerNum] {
			actions = append(actions, fmt.Sprintf(
				"{{b}}raise ##{{_b}} (where ## is at least $%d)",
				max(g.MinimumBet, g.LargestRaise)))
		}
		actions = append(actions, "{{b}}allin{{_b}}")
		buf.WriteString(
			"It's your turn, you can: " + strings.Join(actions, ", ") + "\n\n")
	}
	// All players table
	playersTable := [][]string{
		[]string{
			"{{b}}Players{{_b}}",
			"{{b}}Cash{{_b}}",
			"{{b}}Bet{{_b}}",
			"{{b}}Cards{{_b}}",
		},
	}
	for tablePlayerNum, _ := range g.Players {
		playerRow := []string{g.RenderPlayerName(tablePlayerNum)}
		if tablePlayerNum == g.CurrentDealer {
			playerRow[0] += " (D)"
		}
		if g.PlayerMoney[tablePlayerNum] == 0 && g.Bets[tablePlayerNum] == 0 {
			playerRow = append(playerRow, `{{c "gray"}}Out{{_c}}`)
		} else {
			var cards []string
			if g.FoldedPlayers[tablePlayerNum] {
				cards = []string{`{{c "gray"}}Folded{{_c}}`}
			} else if g.CanSeeHand(playerNum, tablePlayerNum) {
				cards = RenderCards(g.PlayerHands[tablePlayerNum])
			} else {
				cards = []string{
					card.RenderStandard52HiddenFixedWidth(),
					card.RenderStandard52HiddenFixedWidth(),
				}
			}
			playerRow = append(playerRow,
				RenderCash(g.PlayerMoney[tablePlayerNum]),
				RenderCash(g.Bets[tablePlayerNum]), strings.Join(cards, " "))
		}
		playersTable = append(playersTable, playerRow)
	}
	table, err := render.Table(playersTable, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(table)
	g.Log = g.Log.MarkReadFor(player)
	return buf.String(), nil
}

func RenderCards(deck card.Deck) (output []string) {
	for _, c := range deck {
		output = append(output, "{{b}}"+
			c.(card.SuitRankCard).RenderStandard52FixedWidth()+"{{_b}}")
	}
	return
}

func (g *Game) CanSeeHand(playerNum, target int) bool {
	return playerNum == target
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.RemainingPlayers()) < 2
}

func (g *Game) Winners() []string {
	remainingPlayers := g.RemainingPlayers()
	if len(remainingPlayers) == 1 {
		for _, p := range remainingPlayers {
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
