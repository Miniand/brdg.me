package jaipur

import (
	"errors"
	"fmt"

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

func (g *Game) Commands() []command.Command {
	return []command.Command{
		TakeCommand{},
		SellCommand{},
	}
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
	return g.RoundWins[0]+g.RoundWins[1] >= 3
}

func (g *Game) Winners() []string {
	return []string{}
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
	camelWinner := -1
	if g.Camels[0] > g.Camels[1] {
		camelWinner = 0
	} else if g.Camels[0] > g.Camels[1] {
		camelWinner = 1
	}
	if camelWinner != -1 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s won the 5 point camel bonus for having {{b}}%d %s{{_b}}, %s had {{b}}%d{{_b}}",
			g.RenderName(camelWinner),
			g.Camels[camelWinner],
			render.Colour(helper.Plural(g.Camels[camelWinner], "camel"), GoodColours[GoodCamel]),
			g.RenderName(g.Opponent(camelWinner)),
			g.Camels[g.Opponent(camelWinner)],
		)))
		g.Tokens[camelWinner] = append(g.Tokens[camelWinner], CamelBonusPoints)
	}

	winner := -1
	p0Score := helper.IntSum(g.Tokens[0])
	p1Score := helper.IntSum(g.Tokens[1])
	if p0Score > p1Score {
		winner = 0
	} else if p1Score > p0Score {
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
		g.RoundWins[winner]++
	}
	if !g.IsFinished() {
		g.StartRound()
	}
}
