package jaipur

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players []string
	Log     *log.Log

	CurrentPlayer int
	RoundWins     [2]int

	Deck           []int
	Hands, Tokens  [2][]int
	Camels         [2]int
	Bonuses, Goods map[int][]int
	Market         []int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
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
	g.Deck = helper.IntShuffle(Deck())
	g.Market, g.Deck = append([]int{
		GoodCamel,
		GoodCamel,
		GoodCamel,
	}, g.Deck[:2]...), g.Deck[2:]

	g.Camels = [2]int{}
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
	for _, c := range cards {
		switch c {
		case GoodCamel:
			g.Camels[player]++
		default:
			g.Hands[player] = append(g.Hands[player], c)
		}
	}
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
