package for_sale

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
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
	if len(g.ChequeDeck) >= 18 {
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
	}
}

func (g *Game) StartBuyingRound() {
	g.OpenCards, g.BuildingDeck = g.BuildingDeck.PopN(len(g.Players))
	g.OpenCards = g.OpenCards.Sort()
	g.ClearBids()
}

func (g *Game) StartSellingRound() {
	g.OpenCards, g.ChequeDeck = g.ChequeDeck.PopN(len(g.Players))
	g.OpenCards = g.OpenCards.Sort()
	g.ClearBids()
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
	cells := [][]string{}
	cells = append(cells, []string{
		`Your chips:`,
		fmt.Sprintf(`{{b}}%d{{_b}}`, g.Chips[p]),
	})
	buildingRow := []string{`Your buildings:`}
	for _, b := range g.Hands[p] {
		buildingRow = append(buildingRow,
			RenderBuilding(b.(card.SuitRankCard).Rank))
	}
	cells = append(cells, buildingRow)
	chequeRow := []string{`Your cheques:`}
	for _, c := range g.Cheques[p] {
		chequeRow = append(chequeRow,
			RenderCheque(c.(card.SuitRankCard).Rank))
	}
	cells = append(cells, chequeRow)
	table, err := render.Table(cells, 0, 1)
	if err != nil {
		return "", err
	}
	output.WriteString(table)
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
	g.NextBidder()
	return nil
}

func (g *Game) Pass(player int) error {
	if !g.CanBid(player) {
		return errors.New("you are not able to pass at the moment")
	}
	g.TakeFirstOpenCard(player)
	halfBid := g.Bids[player] / 2
	g.Chips[player] -= halfBid
	g.FinishedBidding[player] = true
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
			cheque, g.OpenCards = g.OpenCards.Shift()
			src := c.(card.SuitRankCard)
			g.Cheques[src.Rank] = g.Cheques[src.Rank].Push(cheque)
		}
		g.StartRound()
	}
	return nil
}

func (g *Game) TakeFirstOpenCard(player int) {
	var c card.Card
	c, g.OpenCards = g.OpenCards.Shift()
	g.Hands[player] = g.Hands[player].Push(c).Sort()
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
		g.TakeFirstOpenCard(player)
		g.Chips[player] -= amount
		g.BiddingPlayer = player
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
	for p, b := range g.Bids {
		if b > amount {
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
