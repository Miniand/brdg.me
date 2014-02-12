package modern_art

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"strings"
)

const (
	INITIAL_MONEY = 100
)
const (
	STATE_PLAY_CARD = iota
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
	SUIT_LITE_METAL:   "lm",
	SUIT_YOKO:         "yo",
	SUIT_CHRISTINE_P:  "cp",
	SUIT_KARL_GLITTER: "kg",
	SUIT_KRYPTO:       "kr",
}

var suitColours = map[int]string{
	SUIT_LITE_METAL:   "yellow",
	SUIT_YOKO:         "green",
	SUIT_CHRISTINE_P:  "red",
	SUIT_KARL_GLITTER: "blue",
	SUIT_KRYPTO:       "magenta",
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
	RANK_OPEN:        "op",
	RANK_FIXED_PRICE: "fp",
	RANK_SEALED:      "sl",
	RANK_DOUBLE:      "db",
	RANK_ONCE_AROUND: "oa",
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
	ValueBoard          []map[int]int
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
	pNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	output := bytes.Buffer{}
	// Auction specific
	if g.IsAuction() {
		output.WriteString(fmt.Sprintf("{{b}}Currently auctioning:{{_b}} %s\n",
			RenderCardNames(g.CurrentlyAuctioning)))
		if g.AuctionType() != RANK_SEALED {
			bidder, bid := g.HighestBidder()
			output.WriteString(fmt.Sprintf("{{b}}Current bid:{{_b}} %s by %s\n",
				RenderMoney(bid), render.PlayerName(bidder, g.Players[bidder])))
		}
		output.WriteString("\n")
	}
	// Player money
	output.WriteString(fmt.Sprintf("{{b}}Your money:{{_b}} %s\n\n",
		RenderMoney(g.PlayerMoney[pNum])))
	// Player cards
	output.WriteString("{{b}}Your cards:{{_b}}\n")
	for _, c := range g.PlayerHands[pNum] {
		output.WriteString(fmt.Sprintf("%s\n",
			RenderCardNameCode(c.(card.SuitRankCard))))
	}
	output.WriteString("\n")
	// Players
	cells := [][]string{
		[]string{
			"{{b}}Players{{_b}}",
			"{{b}}Purchases{{_b}}",
		},
	}
	for opNum, oPlayer := range g.Players {
		cards := []string{}
		if len(g.PlayerPurchases[opNum]) > 0 {
			for _, c := range g.PlayerPurchases[opNum] {
				src := c.(card.SuitRankCard)
				cards = append(cards, RenderCardCode(src))
			}
		} else {
			cards = append(cards, `{{c "gray"}}None{{_c}}`)
		}
		cells = append(cells, []string{
			render.PlayerName(opNum, oPlayer),
			strings.Join(cards, " "),
		})
	}
	table, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString(table)
	output.WriteString("\n\n")
	// Artists
	cells = [][]string{
		[]string{
			"{{b}}Artist{{_b}}",
			"{{b}}R1{{_b}}",
			"{{b}}R2{{_b}}",
			"{{b}}R3{{_b}}",
			"{{b}}R4{{_b}}",
			"{{b}}Total{{_b}}",
		},
	}
	for _, s := range suits {
		row := []string{
			RenderSuit(s),
		}
		for i := 0; i < 4; i++ {
			if len(g.ValueBoard) > i {
				row = append(row, RenderMoney(g.ValueBoard[i][s]))
			} else {
				row = append(row, ".")
			}
		}
		row = append(row, RenderMoney(g.SuitValue(s)))
		cells = append(cells, row)
	}
	table, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString(table)
	return output.String(), nil
}

func (g *Game) SuitCardsOnTable(suit int) int {
	count := 0
	for pNum, _ := range g.Players {
		for _, c := range g.PlayerPurchases[pNum] {
			if c.(card.SuitRankCard).Suit == suit {
				count += 1
			}
		}
	}
	for _, c := range g.CurrentlyAuctioning {
		if c.(card.SuitRankCard).Suit == suit {
			count += 1
		}
	}
	return count
}

func (g *Game) SuitValue(suit int) int {
	value := 0
	for _, values := range g.ValueBoard {
		value += values[suit]
	}
	return value
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
	g.ValueBoard = []map[int]int{}
	g.CurrentlyAuctioning = card.Deck{}
	g.Deck = Deck().Shuffle()
	g.Log = &log.Log{}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.State = STATE_PLAY_CARD
	numCards := roundCards[len(g.Players)][g.Round]
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("Start of round %d", g.Round+1)))
	if numCards > 0 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Dealing %d cards to each player", numCards)))
	}
	for i, _ := range g.Players {
		g.PlayerPurchases = map[int]card.Deck{}
		if numCards > 0 {
			cards, remaining := g.Deck.PopN(numCards)
			g.PlayerHands[i] = g.PlayerHands[i].PushMany(cards).Sort()
			g.Deck = remaining
		}
	}
}

func (g *Game) EndRound() {
	g.Log.Add(log.NewPublicMessage("{{b}}It is the end of the round{{_b}}"))
	// Add values to artists
	g.CurrentlyAuctioning = card.Deck{}
	values := map[int]int{}
	scored := map[int]bool{}
	counts := map[int]int{}
	for _, s := range suits {
		counts[s] = g.SuitCardsOnTable(s)
	}
	for _, v := range []int{30, 20, 10} {
		highest := -1
		highestCount := -1
		for _, s := range suits {
			if !scored[s] && counts[s] > highestCount {
				highest = s
				highestCount = counts[s]
			}
		}
		scored[highest] = true
		values[highest] = v
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Adding %s to the value of %s (%d cards)",
			RenderMoney(v), RenderSuit(highest), highestCount)))
	}
	g.ValueBoard = append(g.ValueBoard, values)
	// Pay out purchased cards
	for pNum, p := range g.Players {
		pTotal := 0
		for _, c := range g.PlayerPurchases[pNum] {
			pTotal += g.SuitValue(c.(card.SuitRankCard).Suit)
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Paying %s %s for selling all their cards",
			render.PlayerName(pNum, p), RenderMoney(pTotal))))
		g.PlayerMoney[pNum] += pTotal
	}
	if g.Round == 3 {
		moneyTable := [][]string{}
		for pNum, p := range g.Players {
			moneyTable = append(moneyTable, []string{
				render.PlayerName(pNum, p),
				RenderMoney(g.PlayerMoney[pNum]),
			})
		}
		table, _ := render.Table(moneyTable, 0, 1)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"{{b}}End of the game, final player money:{{_b}}\n%s", table)))
		g.Finished = true
	} else {
		g.Round += 1
		g.NextPlayer()
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
		highestMoney := -1
		highest := []string{}
		for pNum, p := range g.Players {
			if g.PlayerMoney[pNum] > highestMoney {
				highestMoney = g.PlayerMoney[pNum]
				highest = []string{}
			}
			if g.PlayerMoney[pNum] == highestMoney {
				highest = append(highest, p)
			}
		}
		return highest
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
		case RANK_FIXED_PRICE, RANK_DOUBLE:
			// Find the first person without a bid including current player
			for i := 0; i < len(g.Players); i++ {
				p := (i + g.CurrentPlayer) % len(g.Players)
				if _, ok := g.Bids[p]; !ok {
					return []string{g.Players[p]}
				}
			}
		case RANK_ONCE_AROUND:
			// Find the first person without a bid after current player
			highestBid := 0
			for i := 0; i < len(g.Players); i++ {
				p := (1 + i + g.CurrentPlayer) % len(g.Players)
				bid, ok := g.Bids[p]
				if ok && bid > highestBid {
					highestBid = bid
				}
				if !ok && !(highestBid == 0 && p == g.CurrentPlayer) {
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
		case RANK_OPEN, RANK_SEALED, RANK_DOUBLE, RANK_ONCE_AROUND:
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
		case RANK_OPEN, RANK_SEALED, RANK_ONCE_AROUND:
			return g.IsPlayersTurnStr(player)
		}
	}
	return false
}

func (g *Game) CanAdd(player string) bool {
	return g.IsAuction() && g.AuctionType() == RANK_DOUBLE &&
		g.IsPlayersTurnStr(player)
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
	return g.AuctionCard().Rank
}

func (g *Game) AuctionCard() card.SuitRankCard {
	if len(g.CurrentlyAuctioning) == 0 {
		return card.SuitRankCard{}
	}
	return g.CurrentlyAuctioning[len(g.CurrentlyAuctioning)-1].(card.SuitRankCard)
}

func RenderCardNameCode(c card.SuitRankCard) string {
	return RenderInSuit(c.Suit, fmt.Sprintf("(%s) %s",
		CardCode(c), CardName(c)))
}

func RenderCardName(c card.SuitRankCard) string {
	return RenderInSuit(c.Suit, CardName(c))
}

func RenderCardCode(c card.SuitRankCard) string {
	return RenderInSuit(c.Suit, CardCode(c))
}

func RenderSuit(suit int) string {
	return RenderInSuit(suit, suitNames[suit])
}

func RenderInSuit(suit int, s string) string {
	return fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`,
		suitColours[suit], s)
}

func CardName(c card.SuitRankCard) string {
	return fmt.Sprintf("%s - %s", suitNames[c.Suit], rankNames[c.Rank])
}

func CardCode(c card.SuitRankCard) string {
	return fmt.Sprintf("%s%s", suitCodes[c.Suit], rankCodes[c.Rank])
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
	g.CurrentlyAuctioning = card.Deck{}
	return g.AddCardToAuction(player, c)
}

func RenderMoney(amount int) string {
	return fmt.Sprintf(`{{b}}{{c "green"}}$%d{{_c}}{{_b}}`, amount)
}

func (g *Game) AddCardToAuction(player int, c card.SuitRankCard) error {
	remaining, removed := g.PlayerHands[player].Remove(c, 1)
	if removed != 1 {
		return errors.New("You do not have that card in your hand")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s played %s",
		render.PlayerName(player, g.Players[player]),
		RenderCardName(c))))
	g.PlayerHands[player] = remaining
	g.CurrentlyAuctioning = g.CurrentlyAuctioning.Push(c)
	g.Bids = map[int]int{}
	g.State = STATE_AUCTION
	if g.SuitCardsOnTable(c.Suit) >= 5 {
		g.EndRound()
	}
	return nil
}

func RenderCardNames(d card.Deck) string {
	cardStrings := []string{}
	for _, c := range d {
		cardStrings = append(cardStrings,
			RenderCardName(c.(card.SuitRankCard)))
	}
	return strings.Join(cardStrings, " and ")
}

func (g *Game) SettleAuction(winner, price int) {
	g.PlayerMoney[winner] -= price
	g.PlayerPurchases[winner] = g.PlayerPurchases[winner].
		PushMany(g.CurrentlyAuctioning).Sort()
	paidTo := "the bank"
	if winner != g.CurrentPlayer {
		g.PlayerMoney[g.CurrentPlayer] += price
		paidTo = render.PlayerName(g.CurrentPlayer, g.Players[g.CurrentPlayer])
	}
	g.State = STATE_PLAY_CARD
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s bought %s, paying %s to %s",
		render.PlayerName(winner, g.Players[winner]),
		RenderCardNames(g.CurrentlyAuctioning), RenderMoney(price), paidTo)))
	g.NextPlayer()
}

func (g *Game) Pass(player int) error {
	if !g.CanPass(g.Players[player]) {
		return errors.New("You're not able to pass at the moment")
	}
	g.Bids[player] = 0
	if g.AuctionType() != RANK_SEALED {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s passed",
			render.PlayerName(player, g.Players[player]))))
	}
	switch g.AuctionType() {
	case RANK_FIXED_PRICE:
		if len(g.Bids) == len(g.Players) {
			g.SettleAuction(g.CurrentPlayer, g.Bids[g.CurrentPlayer])
		}
	default:
		if len(g.WhoseTurn()) == 0 {
			g.SettleAuction(g.HighestBidder())
		}
	}
	return nil
}

func (g *Game) Bid(player, amount int) error {
	if !g.CanBid(g.Players[player]) {
		return errors.New("You're not able to bid at the moment")
	}
	if amount > g.PlayerMoney[player] {
		return fmt.Errorf(
			"You must not bid higher than the money you have, which is $%d",
			g.PlayerMoney[player])
	}
	if g.AuctionType() != RANK_SEALED {
		_, highestBid := g.HighestBidder()
		if amount <= highestBid {
			return fmt.Errorf("You must bid higher than $%d", highestBid)
		}
	}
	g.Bids[player] = amount
	if g.AuctionType() != RANK_SEALED {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s bid %s",
			render.PlayerName(player, g.Players[player]),
			RenderMoney(amount))))
	}
	if len(g.WhoseTurn()) == 0 {
		g.SettleAuction(g.HighestBidder())
	}
	return nil
}

func (g *Game) AddCard(player int, c card.SuitRankCard) error {
	if !g.CanAdd(g.Players[player]) {
		return errors.New("You're not able to add a card at the moment")
	}
	if g.AuctionCard().Suit != c.Suit {
		return errors.New("The artist of the card must match the existing one")
	}
	if c.Rank == RANK_DOUBLE {
		return errors.New("You are not allowed to add a second double auction")
	}
	err := g.AddCardToAuction(player, c)
	if err != nil {
		return err
	}
	g.CurrentPlayer = player
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
