package jaipur

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players []string
	Log     *log.Log

	CurrentPlayer int
	RoundWins     [2]int

	Deck                            []int
	Hands, Tokens                   [2][]int
	Camels, BonusTokens, GoodTokens [2]int
	Bonuses, Goods                  map[int][]int
	Market                          []int
}

func (g *Game) Commands(player string) []command.Command {
	commands := []command.Command{}
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return commands
	}
	if g.CanTake(pNum) {
		commands = append(commands, TakeCommand{})
	}
	if g.CanSell(pNum) {
		commands = append(commands, SellCommand{})
	}
	return commands
}

func (g *Game) Name() string {
	return "Jaipur"
}

func (g *Game) Identifier() string {
	return "jaipur"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("only two players allowed")
	}
	g.Players = players
	g.Log = log.New()

	g.RoundWins = [2]int{}

	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"It is the start of the round, starting market with {{b}}3 %s{{_b}}",
		render.Colour("camels", GoodColours[GoodCamel]),
	)))
	g.Deck = helper.IntShuffle(Deck())
	g.Market = []int{GoodCamel, GoodCamel, GoodCamel}
	g.ReplenishMarket()

	g.Camels = [2]int{}
	g.BonusTokens = [2]int{}
	g.GoodTokens = [2]int{}
	g.Hands = [2][]int{}
	g.Tokens = [2][]int{}
	for p := range g.Players {
		var hand []int
		hand, g.Deck = g.Deck[:5], g.Deck[5:]
		g.ReceiveCards(p, hand)
		g.Tokens[p] = []int{}
	}

	g.Goods = map[int][]int{}
	for _, good := range TradeGoods {
		g.Goods[good] = append([]int{}, TradeGoodTokens[good]...)
	}

	g.Bonuses = map[int][]int{}
	for i := MinTradeBonus; i <= MaxTradeBonus; i++ {
		g.Bonuses[i] = helper.IntShuffle(TradeBonuses[i])
	}
}

func (g *Game) ReceiveCards(player int, cards []int) {
	var drawnGoods, drawnCamels int
	for _, c := range cards {
		switch c {
		case GoodCamel:
			g.Camels[player]++
			drawnCamels++
		default:
			g.Hands[player] = append(g.Hands[player], c)
			drawnGoods++
		}
	}
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You drew %s",
		render.CommaList(RenderGoods(helper.IntSort(cards))),
	), []string{g.Players[player]}))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"%s drew {{b}}%d %s{{_b}} and {{b}}%d %s{{_b}}",
		g.RenderName(player),
		drawnGoods,
		helper.Plural(drawnGoods, "good"),
		drawnCamels,
		render.Colour(helper.Plural(drawnCamels, "camel"), GoodColours[GoodCamel]),
	), []string{g.Players[g.Opponent(player)]}))
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.RoundWins[0] == 2 || g.RoundWins[1] == 2
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	if g.RoundWins[0] == 2 {
		return []string{g.Players[0]}
	}
	return []string{g.Players[1]}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 2
}

func (g *Game) ReplenishMarket() bool {
	n := 5 - len(g.Market)
	if n == 0 {
		return true
	} else if len(g.Deck) < n {
		g.EndRound()
		return false
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Drew %s from the deck and added %s to the market",
		render.CommaList(RenderGoods(helper.IntSort(g.Deck[:n]))),
		helper.Plural(n, "it"),
	)))
	g.Market = append(g.Market, g.Deck[:n]...)
	g.Deck = g.Deck[n:]
	return true
}

func (g *Game) Opponent(player int) int {
	return (player + 1) % 2
}

func (g *Game) EndRound() {
	logBuf := bytes.Buffer{}
	camelWinner := -1
	if g.Camels[0] > g.Camels[1] {
		camelWinner = 0
	} else if g.Camels[0] > g.Camels[1] {
		camelWinner = 1
	}
	if camelWinner != -1 {
		logBuf.WriteString(fmt.Sprintf(
			"%s won the 5 point camel bonus for having {{b}}%d %s{{_b}}, %s had {{b}}%d{{_b}}\n",
			g.RenderName(camelWinner),
			g.Camels[camelWinner],
			render.Colour(helper.Plural(g.Camels[camelWinner], "camel"), GoodColours[GoodCamel]),
			g.RenderName(g.Opponent(camelWinner)),
			g.Camels[g.Opponent(camelWinner)],
		))
		g.Tokens[camelWinner] = append(g.Tokens[camelWinner], CamelBonusPoints)
		g.BonusTokens[camelWinner]++
	}

	scores := map[int]int{}

	for p := range g.Players {
		scores[p] = helper.IntSum(g.Tokens[p])
		logBuf.WriteString(fmt.Sprintf(
			"%s had {{b}}%d{{_b}} %s from {{b}}%d{{_b}} bonus %s and {{b}}%d{{_b}} good %s\n",
			g.RenderName(p),
			scores[p],
			helper.Plural(scores[p], "point"),
			g.BonusTokens[p],
			helper.Plural(g.BonusTokens[p], "token"),
			g.GoodTokens[p],
			helper.Plural(g.GoodTokens[p], "token"),
		))
	}

	winner := -1
	if scores[0] > scores[1] {
		winner = 0
	} else if scores[1] > scores[0] {
		winner = 1
	} else if g.BonusTokens[0] > g.BonusTokens[1] {
		winner = 0
	} else if g.BonusTokens[1] > g.BonusTokens[0] {
		winner = 1
	} else if g.GoodTokens[0] > g.GoodTokens[1] {
		winner = 0
	} else if g.GoodTokens[1] > g.GoodTokens[0] {
		winner = 1
	}
	if winner != -1 {
		logBuf.WriteString(fmt.Sprintf(
			"%s won the round",
			g.RenderName(winner),
		))
		g.RoundWins[winner]++
	} else {
		logBuf.WriteString("Against all odds, the round was tied and will be replayed")
	}
	g.Log.Add(log.NewPublicMessage(logBuf.String()))
	if !g.IsFinished() {
		g.StartRound()
	}
}

func ParseGoodStr(input string) ([]int, error) {
	words := strings.Split(input, " ")
	goods := []int{}
	quantity := -1
	for _, w := range words {
		// See if it's a quantity first.
		if quantity == -1 {
			q, err := strconv.Atoi(w)
			if err == nil {
				if q <= 0 {
					return nil, errors.New("quantities must be 1 or larger")
				}
				quantity = q
				continue
			}
		}
		good, err := helper.MatchStringInStringMap(w, GoodStringsPl)
		if err != nil {
			return nil, err
		}
		if quantity == -1 {
			quantity = 1
		}
		for i := 0; i < quantity; i++ {
			goods = append(goods, good)
		}
		quantity = -1
	}
	return goods, nil
}
