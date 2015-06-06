package for_sale

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	BuyingPhase  = 0
	SellingPhase = 1
	GameFinished = 2
)

type Game struct {
	Players         []string
	BuildingDeck    card.Deck
	ChequeDeck      card.Deck
	OpenCards       card.Deck
	Hands           map[int]card.Deck
	Cheques         map[int]card.Deck
	Chips           map[int]int
	BiddingPlayer   int
	Bids            map[int]int
	FinishedBidding map[int]bool
	Log             *log.Log
}

func (g *Game) Name() string {
	return "For Sale"
}

func (g *Game) Identifier() string {
	return "for_sale"
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		BidCommand{},
		PassCommand{},
		PlayCommand{},
	}
}

func (g *Game) Start(players []string) error {
	if len(players) < 3 || len(players) > 5 {
		return errors.New("must have between 3 and 5 players")
	}
	g.Log = log.New()
	g.Players = players
	g.BuildingDeck = BuildingDeck().Shuffle()
	g.ChequeDeck = ChequeDeck().Shuffle()
	g.Hands = map[int]card.Deck{}
	g.Cheques = map[int]card.Deck{}
	g.Chips = map[int]int{}
	g.Bids = map[int]int{}
	g.FinishedBidding = map[int]bool{}
	for p, _ := range g.Players {
		g.Hands[p] = card.Deck{}
		g.Cheques[p] = card.Deck{}
		g.Chips[p] = 15
		g.Bids[p] = 0
		g.FinishedBidding[p] = false
	}
	if len(players) == 3 {
		g.Log.Add(log.NewPublicMessage(
			"Removing two building and cheque cards for 3 player game"))
		_, g.BuildingDeck = g.BuildingDeck.PopN(2)
		_, g.ChequeDeck = g.ChequeDeck.PopN(2)
	}
	g.StartRound()
	return nil
}

func (g *Game) CurrentPhase() int {
	if len(g.BuildingDeck) > 0 ||
		(len(g.OpenCards) > 0 && len(g.ChequeDeck) >= 18) {
		return BuyingPhase
	} else if len(g.ChequeDeck) > 0 || len(g.OpenCards) > 0 {
		return SellingPhase
	}
	return GameFinished
}

func (g *Game) StartRound() {
	switch g.CurrentPhase() {
	case BuyingPhase:
		g.StartBuyingRound()
	case SellingPhase:
		g.StartSellingRound()
	case GameFinished:
		output := bytes.NewBufferString(
			"{{b}}The game has finished!  The scores are:{{_b}}\n")
		playerScores := [][]interface{}{}
		for pNum, p := range g.Players {
			playerScores = append(playerScores, []interface{}{
				render.PlayerName(pNum, p),
				fmt.Sprintf("{{b}}%d{{_b}}", g.DeckValue(g.Cheques[pNum])),
			})
		}
		table := render.Table(playerScores, 0, 1)
		output.WriteString(table)
		g.Log.Add(log.NewPublicMessage(output.String()))
	}
}

func (g *Game) StartBuyingRound() {
	g.OpenCards, g.BuildingDeck = g.BuildingDeck.PopN(len(g.Players))
	g.OpenCards = g.OpenCards.Sort()
	g.ClearBids()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(`Drew new buildings: %s`,
		strings.Join(RenderCards(g.OpenCards, RenderBuilding), " "))))
}

func (g *Game) StartSellingRound() {
	g.OpenCards, g.ChequeDeck = g.ChequeDeck.PopN(len(g.Players))
	g.OpenCards = g.OpenCards.Sort()
	g.ClearBids()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(`Drew new cheques: %s`,
		strings.Join(RenderCards(g.OpenCards, RenderCheque), " "))))
	if g.Hands[0].Len() == 1 {
		// Autoplay the final card
		fmt.Println("autoplaying")
		for p, _ := range g.Players {
			g.Play(p, g.Hands[p][0].(card.SuitRankCard).Rank)
		}
	}
}

func (g *Game) ClearBids() {
	for p, _ := range g.Players {
		g.Bids[p] = 0
		g.FinishedBidding[p] = false
	}
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	output := bytes.NewBuffer([]byte{})
	switch g.CurrentPhase() {
	case BuyingPhase:
		output.WriteString(fmt.Sprintf("Buildings available: %s\n",
			strings.Join(RenderCards(g.OpenCards, RenderBuilding), " ")))
		currentBidText := `{{c "gray"}}none{{_c}}`
		if highestPlayer, highestAmount := g.HighestBid(); highestAmount > 0 {
			currentBidText = fmt.Sprintf(`{{b}}%d{{_b}} by %s`, highestAmount,
				render.PlayerName(highestPlayer, g.Players[highestPlayer]))
		}
		output.WriteString(fmt.Sprintf("Current bid: %s\n", currentBidText))
		output.WriteString(fmt.Sprintf("Your bid: {{b}}%d{{_b}}\n",
			g.Bids[p]))
		remainingPlayers := []string{}
		for remP, remPName := range g.Players {
			if !g.FinishedBidding[remP] {
				remainingPlayers = append(remainingPlayers,
					render.PlayerName(remP, remPName))
			}
		}
		output.WriteString(fmt.Sprintf("Remaining players: %s\n\n",
			render.CommaList(remainingPlayers)))
	case SellingPhase:
		output.WriteString(fmt.Sprintf("Cheques available: %s\n\n",
			strings.Join(RenderCards(g.OpenCards, RenderCheque), " ")))
	}
	output.WriteString(fmt.Sprintf("Your chips: {{b}}%d{{_b}}\n", g.Chips[p]))
	output.WriteString(fmt.Sprintf("Your buildings: %s\n",
		strings.Join(RenderCards(g.Hands[p], RenderBuilding), " ")))
	output.WriteString(fmt.Sprintf("Your cheques: %s",
		strings.Join(RenderCards(g.Cheques[p], RenderCheque), " ")))

	if !g.IsFinished() {
		var (
			rounds    int
			roundType string
		)
		switch g.CurrentPhase() {
		case BuyingPhase:
			rounds = (g.BuildingDeck.Len() / len(g.Players)) + 1
			roundType = "buying"
		case SellingPhase:
			rounds = (g.ChequeDeck.Len() / len(g.Players)) + 1
			roundType = "selling"
		}
		output.WriteString(fmt.Sprintf(
			"\n\n{{b}}%d{{_b}} %s %s remaining",
			rounds,
			roundType,
			helper.Plural(rounds, "round"),
		))
	}
	return output.String(), nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) DeckValue(deck card.Deck) int {
	value := 0
	for _, c := range deck {
		value += c.(card.SuitRankCard).Rank
	}
	return value
}

func (g *Game) IsFinished() bool {
	return len(g.OpenCards) == 0 && len(g.BuildingDeck) == 0 &&
		len(g.ChequeDeck) == 0
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	players := []string{}
	maxMoney := -1
	maxChips := -1
	for pn, p := range g.Players {
		pMoney := g.DeckValue(g.Cheques[pn])
		if pMoney > maxMoney ||
			(pMoney == maxMoney && g.Chips[pn] > maxChips) {
			players = []string{}
			maxMoney = pMoney
			maxChips = g.Chips[pn]
		}
		if pMoney == maxMoney && g.Chips[pn] == maxChips {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	} else if g.CurrentPhase() == BuyingPhase {
		return g.WhoseTurnBuying()
	}
	return g.WhoseTurnSelling()
}

func (g *Game) WhoseTurnBuying() []string {
	return []string{g.Players[g.BiddingPlayer]}
}

func (g *Game) WhoseTurnSelling() []string {
	players := []string{}
	for pn, p := range g.Players {
		if !g.FinishedBidding[pn] {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) GameLog() *log.Log {
	return g.Log
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

func (g *Game) ParsePlayer(player string) (int, error) {
	for pn, p := range g.Players {
		if p == player {
			return pn, nil
		}
	}
	return 0, fmt.Errorf("could not find player %s", player)
}

func (g *Game) CanBid(player int) bool {
	return !g.IsFinished() && g.CurrentPhase() == BuyingPhase &&
		g.BiddingPlayer == player
}

func (g *Game) Bid(player, amount int) error {
	if !g.CanBid(player) {
		return errors.New("you are not able to bid at the moment")
	}
	if amount > g.Chips[player] {
		return fmt.Errorf("cannot bid %d, you only have %d",
			amount, g.Chips[player])
	}
	if _, highest := g.HighestBid(); amount <= highest {
		return fmt.Errorf("you must bid higher than %d", highest)
	}
	g.Bids[player] = amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s bid {{b}}%d{{_b}}",
		render.PlayerName(player, g.Players[player]), amount)))
	g.NextBidder()
	return nil
}

func (g *Game) Pass(player int) error {
	if !g.CanBid(player) {
		return errors.New("you are not able to pass at the moment")
	}
	c := g.TakeFirstOpenCard(player)
	halfBid := g.Bids[player] / 2
	g.Chips[player] -= halfBid
	g.FinishedBidding[player] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s passed, paying {{b}}%d{{_b}} for %s",
		render.PlayerName(player, g.Players[player]), halfBid,
		RenderBuilding(c.Rank))))
	g.NextBidder()
	return nil
}

func (g *Game) CanPlay(player int) bool {
	return !g.IsFinished() && g.CurrentPhase() == SellingPhase &&
		!g.FinishedBidding[player]
}

func (g *Game) Play(player, building int) error {
	var cheque card.Card
	if !g.CanPlay(player) {
		return errors.New("you are not able to play a building card at the moment")
	}
	remaining, n := g.Hands[player].Remove(card.SuitRankCard{
		Rank: building,
	}, 1)
	if n == 0 {
		return errors.New("you don't have that card in your hand")
	}
	g.Hands[player] = remaining
	g.Bids[player] = building
	g.FinishedBidding[player] = true
	if len(g.WhoseTurn()) == 0 {
		played := card.Deck{}
		for p, b := range g.Bids {
			played = append(played, card.SuitRankCard{
				Suit: b,
				Rank: p,
			})
		}
		for _, c := range played.Sort() {
			src := c.(card.SuitRankCard)
			p := src.Rank
			cheque, g.OpenCards = g.OpenCards.Shift()
			g.Cheques[p] = g.Cheques[p].Push(cheque)
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`%s sold %s for %s`, render.PlayerName(p, g.Players[p]),
				RenderBuilding(src.Suit),
				RenderCheque(cheque.(card.SuitRankCard).Rank))))
		}
		g.StartRound()
	}
	return nil
}

func (g *Game) TakeFirstOpenCard(player int) card.SuitRankCard {
	var c card.Card
	c, g.OpenCards = g.OpenCards.Shift()
	g.Hands[player] = g.Hands[player].Push(c).Sort()
	return c.(card.SuitRankCard)
}

func (g *Game) NextBidder() {
	remaining := 0
	for _, b := range g.FinishedBidding {
		if !b {
			remaining += 1
		}
	}
	if remaining == 1 {
		// Last remaining player takes the last building for the full price.
		player, amount := g.HighestBid()
		c := g.TakeFirstOpenCard(player)
		g.Chips[player] -= amount
		g.BiddingPlayer = player
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s is the last player, paying {{b}}%d{{_b}} for %s",
			render.PlayerName(player, g.Players[player]), amount,
			RenderBuilding(c.Rank))))
		g.StartRound()
		return
	}
	for {
		g.BiddingPlayer = (g.BiddingPlayer + 1) % len(g.Players)
		if !g.FinishedBidding[g.BiddingPlayer] {
			break
		}
	}
}

func (g *Game) HighestBid() (player, amount int) {
	amount = -1
	for p, b := range g.Bids {
		if !g.FinishedBidding[p] && b > amount {
			player = p
			amount = b
		}
	}
	return
}

func BuildingDeck() card.Deck {
	d := card.Deck{}
	for i := 1; i <= 20; i++ {
		d = d.Push(card.SuitRankCard{
			Rank: i,
		})
	}
	return d
}

func ChequeDeck() card.Deck {
	d := card.Deck{}
	for i := 1; i <= 20; i++ {
		c := card.SuitRankCard{
			Rank: i,
		}
		if i < 3 {
			c.Rank = 0
		}
		d = d.Push(c)
	}
	return d
}

func RenderBuilding(value int) string {
	return fmt.Sprintf(`{{b}}{{c "green"}}%d{{_c}}{{_b}}`, value)
}

func RenderCheque(value int) string {
	return fmt.Sprintf(`{{b}}{{c "blue"}}%d{{_c}}{{_b}}`, value)
}

func RenderCards(deck card.Deck, renderer func(int) string) []string {
	output := []string{}
	for _, c := range deck {
		output = append(output, renderer(c.(card.SuitRankCard).Rank))
	}
	return output
}
