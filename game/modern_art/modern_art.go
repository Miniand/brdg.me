package modern_art

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"strings"
)

const (
	INITIAL_MONEY = 100
)
const (
	STATE_PLAY_CARD = iota
	STATE_ADD_DOUBLE
	STATE_AUCTION
)
const (
	SUIT_LITE_METAL = iota
	SUIT_YOKO
	SUIT_CHRISTINE_P
	SUIT_KARL_GLITTER
	SUIT_KRYPTO
)
const (
	RANK_OPEN = iota
	RANK_FIXED_PRICE
	RANK_SEALED
	RANK_DOUBLE
	RANK_ONCE_AROUND
)

var roundCards = map[int]map[int]int{
	3: map[int]int{
		0: 10,
		1: 6,
		2: 6,
		3: 0,
	},
	4: map[int]int{
		0: 9,
		1: 4,
		2: 4,
		3: 0,
	},
	5: map[int]int{
		0: 8,
		1: 3,
		2: 3,
		3: 0,
	},
}

var suits = []int{
	SUIT_LITE_METAL,
	SUIT_YOKO,
	SUIT_CHRISTINE_P,
	SUIT_KARL_GLITTER,
	SUIT_KRYPTO,
}

var suitNames = map[int]string{
	SUIT_LITE_METAL:   "Lite Metal",
	SUIT_YOKO:         "Yoko",
	SUIT_CHRISTINE_P:  "Christine P",
	SUIT_KARL_GLITTER: "Karl Glitter",
	SUIT_KRYPTO:       "Krypto",
}

var suitCodes = map[int]string{
	SUIT_LITE_METAL:   "LM",
	SUIT_YOKO:         "YO",
	SUIT_CHRISTINE_P:  "CP",
	SUIT_KARL_GLITTER: "KG",
	SUIT_KRYPTO:       "KR",
}

var ranks = []int{
	RANK_OPEN,
	RANK_FIXED_PRICE,
	RANK_SEALED,
	RANK_DOUBLE,
	RANK_ONCE_AROUND,
}

var rankNames = map[int]string{
	RANK_OPEN:        "Open",
	RANK_FIXED_PRICE: "Fixed Price",
	RANK_SEALED:      "Sealed",
	RANK_DOUBLE:      "Double",
	RANK_ONCE_AROUND: "Once Around",
}

var rankCodes = map[int]string{
	RANK_OPEN:        "OP",
	RANK_FIXED_PRICE: "FP",
	RANK_SEALED:      "SL",
	RANK_DOUBLE:      "DB",
	RANK_ONCE_AROUND: "OA",
}

var cardDistribution = map[int]map[int]int{
	SUIT_LITE_METAL: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 2,
		RANK_SEALED:      2,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_YOKO: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 2,
	},
	SUIT_CHRISTINE_P: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      2,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_KARL_GLITTER: map[int]int{
		RANK_OPEN:        3,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      3,
		RANK_ONCE_AROUND: 3,
	},
	SUIT_KRYPTO: map[int]int{
		RANK_OPEN:        4,
		RANK_FIXED_PRICE: 3,
		RANK_SEALED:      3,
		RANK_DOUBLE:      3,
		RANK_ONCE_AROUND: 3,
	},
}

type Game struct {
	Players             []string
	PlayerMoney         map[int]int
	PlayerHands         map[int]card.Deck
	PlayerPurchases     map[int]card.Deck
	State               int
	Round               int
	Deck                card.Deck
	Log                 *log.Log
	CurrentPlayer       int
	ValueBoard          map[int]map[int]int
	Finished            bool
	CurrentlyAuctioning card.Deck
	Bids                map[int]int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
		PriceCommand{},
		AddCommand{},
		BidCommand{},
		BuyCommand{},
		PassCommand{},
	}
}

func (g *Game) Name() string {
	return "Modern Art"
}

func (g *Game) Identifier() string {
	return "modern_art"
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

func (g *Game) Start(players []string) error {
	if len(players) < 3 || len(players) > 5 {
		return errors.New("Modern Art requires between 3 and 5 players")
	}
	g.Players = players
	g.PlayerMoney = map[int]int{}
	g.PlayerHands = map[int]card.Deck{}
	for i, _ := range g.Players {
		g.PlayerMoney[i] = INITIAL_MONEY
		g.PlayerHands[i] = card.Deck{}
	}
	g.CurrentlyAuctioning = card.Deck{}
	g.Deck = Deck().Shuffle()
	g.Log = &log.Log{}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	numCards := roundCards[len(g.Players)][g.Round]
	if numCards <= 0 {
		return
	}
	for i, _ := range g.Players {
		g.PlayerPurchases = map[int]card.Deck{}
		cards, remaining := g.Deck.PopN(numCards)
		g.PlayerHands[i] = g.PlayerHands[i].PushMany(cards)
		g.Deck = remaining
	}
}

func (g *Game) EndRound() {
	if g.Round == 3 {
		g.Finished = true
	} else {
		g.Round += 1
		g.StartRound()
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Finished
}

func (g *Game) Winners() []string {
	if g.IsFinished() {
		return []string{}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	switch g.State {
	case STATE_PLAY_CARD:
		return []string{g.Players[g.CurrentPlayer]}
	case STATE_AUCTION:
		switch g.AuctionType() {
		case RANK_OPEN:
			players := []string{}
			highestBidder, _ := g.HighestBidder()
			for pNum, p := range g.Players {
				if bid, ok := g.Bids[pNum]; pNum != highestBidder &&
					(!ok || bid > 0) {
					players = append(players, p)
				}
			}
			return players
		case RANK_FIXED_PRICE:
			for i := 0; i < len(g.Players); i++ {
				p := (i + g.CurrentPlayer) % len(g.Players)
				if _, ok := g.Bids[p]; !ok {
					return []string{g.Players[p]}
				}
			}
		case RANK_SEALED:
			players := []string{}
			for pNum, p := range g.Players {
				if _, ok := g.Bids[pNum]; !ok {
					players = append(players, p)
				}
			}
			return players
		}
	}
	return []string{}
}

func (g *Game) HighestBidder() (player, bid int) {
	bid = -1
	for i := g.CurrentPlayer; i < g.CurrentPlayer+len(g.Players); i++ {
		p := i % len(g.Players)
		if g.Bids[p] > bid {
			player = p
			bid = g.Bids[p]
		}
	}
	return
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) ParseCardString(s string) (card.SuitRankCard, error) {
	return card.SuitRankCard{}, nil
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
}

func (g *Game) CanPlay(player string) bool {
	return !g.IsFinished() && g.IsPlayersTurnStr(player) &&
		g.State == STATE_PLAY_CARD
}

func (g *Game) CanPass(player string) bool {
	if g.IsAuction() {
		switch g.AuctionType() {
		case RANK_OPEN, RANK_SEALED:
			return g.IsPlayersTurnStr(player)
		case RANK_FIXED_PRICE:
			return player != g.Players[g.CurrentPlayer] &&
				g.IsPlayersTurnStr(player)
		}
	}
	return false
}

func (g *Game) CanBid(player string) bool {
	if g.IsAuction() {
		switch g.AuctionType() {
		case RANK_OPEN, RANK_SEALED:
			return g.IsPlayersTurnStr(player)
		}
	}
	return false
}

func (g *Game) CanAdd(player string) bool {
	return g.IsPlayersTurnStr(player) && g.State == STATE_AUCTION &&
		len(g.CurrentlyAuctioning) == 1 &&
		g.CurrentlyAuctioning[0].(card.SuitRankCard).Rank == RANK_DOUBLE
}

func (g *Game) CanBuy(player string) bool {
	return g.IsAuction() && g.AuctionType() == RANK_FIXED_PRICE &&
		g.IsPlayersTurnStr(player) && g.Players[g.CurrentPlayer] != player
}

func (g *Game) CanSetPrice(player string) bool {
	return g.IsAuction() && g.AuctionType() == RANK_FIXED_PRICE &&
		g.IsPlayersTurnStr(player) && g.Players[g.CurrentPlayer] == player
}

func (g *Game) IsAuction() bool {
	return g.State == STATE_AUCTION
}

func (g *Game) AuctionType() int {
	if !g.IsAuction() || len(g.CurrentlyAuctioning) == 0 {
		return -1
	}
	return g.CurrentlyAuctioning[len(g.CurrentlyAuctioning)-1].(card.SuitRankCard).Rank
}

func (g *Game) SetPrice(player, price int) error {
	if !g.CanSetPrice(g.Players[player]) {
		return errors.New("You're not able to set the price at the moment")
	}
	if price <= 0 {
		return errors.New("The price you set must be higher than 0")
	}
	if price > g.PlayerMoney[player] {
		return errors.New("You can't set the price higher than your current money")
	}
	g.Bids[player] = price
	return nil
}

func (g *Game) Buy(player int) error {
	if !g.CanBuy(g.Players[player]) {
		return errors.New("You're not able to buy the card at the moment")
	}
	price := g.Bids[g.CurrentPlayer]
	if price > g.PlayerMoney[player] {
		return errors.New("You don't have enough money to buy the card")
	}
	g.SettleAuction(player, price)
	return nil
}

func (g *Game) PlayCard(player int, c card.SuitRankCard) error {
	if !g.CanPlay(g.Players[player]) {
		return errors.New("You're not able to play a card at the moment")
	}
	remaining, removed := g.PlayerHands[player].Remove(c, 1)
	if removed != 1 {
		return errors.New("You do not have that card in your hand")
	}
	g.PlayerHands[player] = remaining
	g.CurrentlyAuctioning = card.Deck{c}
	g.Bids = map[int]int{}
	if c.Rank == RANK_DOUBLE {
		g.State = STATE_ADD_DOUBLE
	} else {
		g.State = STATE_AUCTION
	}
	return nil
}

func (g *Game) SettleAuction(winner, price int) {
	g.PlayerMoney[winner] -= price
	g.PlayerPurchases[winner] = g.PlayerPurchases[winner].
		PushMany(g.CurrentlyAuctioning)
	if winner != g.CurrentPlayer {
		g.PlayerMoney[g.CurrentPlayer] += price
	}
	g.State = STATE_PLAY_CARD
	g.NextPlayer()
}

func (g *Game) Pass(player int) error {
	if !g.CanPass(g.Players[player]) {
		return errors.New("You're not able to pass at the moment")
	}
	g.Bids[player] = 0
	switch g.AuctionType() {
	case RANK_OPEN, RANK_SEALED:
		if len(g.WhoseTurn()) == 0 {
			g.SettleAuction(g.HighestBidder())
		}
	case RANK_FIXED_PRICE:
		if len(g.Bids) == len(g.Players) {
			g.SettleAuction(g.CurrentPlayer, g.Bids[g.CurrentPlayer])
		}
	}
	return nil
}

func (g *Game) Bid(player, amount int) error {
	if !g.CanBid(g.Players[player]) {
		return errors.New("You're not able to bid at the moment")
	}
	g.Bids[player] = amount
	switch g.AuctionType() {
	case RANK_SEALED:
		if len(g.WhoseTurn()) == 0 {
			g.SettleAuction(g.HighestBidder())
		}
	}
	return nil
}

func (g *Game) AddCard(player int, c card.SuitRankCard) error {
	if !g.CanAdd(g.Players[player]) {
		return errors.New("You're not able to add a card at the moment")
	}
	return nil
}

func (g *Game) PlayerFromString(s string) (int, error) {
	for i, p := range g.Players {
		if p == s {
			return i, nil
		}
	}
	return 0, errors.New("Could not find player")
}

func (g *Game) IsPlayersTurn(player int) bool {
	for _, p := range g.WhoseTurn() {
		if p == g.Players[player] {
			return true
		}
	}
	return false
}

func (g *Game) IsPlayersTurnStr(player string) bool {
	p, err := g.PlayerFromString(player)
	if err != nil {
		return false
	}
	return g.IsPlayersTurn(p)
}

func Deck() card.Deck {
	d := card.Deck{}
	for suit, suitCards := range cardDistribution {
		for rank, n := range suitCards {
			for i := 0; i < n; i++ {
				d = d.Push(card.SuitRankCard{suit, rank})
			}
		}
	}
	return d
}

func ParseCard(s string) (card.SuitRankCard, error) {
	raw := strings.ToUpper(strings.TrimSpace(s))
	c := card.SuitRankCard{}
	found := false
	for code, prefix := range suitCodes {
		upperPrefix := strings.ToUpper(prefix)
		if strings.HasPrefix(raw, upperPrefix) {
			found = true
			c.Suit = code
			raw = strings.TrimPrefix(raw, upperPrefix)
			break
		}
	}
	if !found {
		return c, errors.New("Could not find the artist in card code")
	}
	for code, suffix := range rankCodes {
		upperSuffix := strings.ToUpper(suffix)
		if strings.HasSuffix(raw, upperSuffix) {
			found = true
			c.Rank = code
			raw = strings.TrimSuffix(raw, upperSuffix)
			break
		}
	}
	if !found {
		return c, errors.New("Could not find the auction type in card code")
	}
	return c, nil
}
