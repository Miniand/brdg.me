package seven_wonders_duel

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players []string
	Log     *log.Log

	Layout Layout

	PlayerCoins [2]int
	PlayerCards [2][]int
}

func (g *Game) Commands(player string) []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "7 Wonders: Duel"
}

func (g *Game) Identifier() string {
	return "7_wonders_duel"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("2 players only")
	}
	g.Players = players
	g.Log = log.New()
	g.PlayerCoins = [2]int{}
	g.PlayerCards = [2][]int{{}, {}}
	g.Layout = Age1Layout()
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) ModifyCoins(player, amount int) {
	if amount == 0 {
		return
	}
	verb := "gained"
	logAmount := amount
	if amount < 0 {
		if g.PlayerCoins[player]-amount < 0 {
			amount = g.PlayerCoins[player]
		}
		verb = "lost"
		logAmount = -amount
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s %s %d %s",
		g.PlayerName(player),
		verb,
		logAmount,
		helper.Plural(logAmount, "coin"),
	)))
	g.PlayerCoins[player] += amount
}

func (g *Game) PlayerCardTypeCount(player, cardType int) int {
	num := 0
	for _, c := range g.PlayerCards[player] {
		if Cards[c].Type == cardType {
			num++
		}
	}
	return num
}

func (g *Game) GreatestCardCount(cardTypes ...int) int {
	num := 0
	for p := range g.Players {
		pNum := 0
		for _, ct := range cardTypes {
			pNum += g.PlayerCardTypeCount(p, ct)
		}
		if pNum > num {
			num = pNum
		}
	}
	return num
}

func Opponent(player int) int {
	return (player + 1) % 2
}
