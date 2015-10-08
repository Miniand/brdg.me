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

	Coins [2]int
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
	g.Coins = [2]int{}
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
		if g.Coins[player]-amount < 0 {
			amount = g.Coins[player]
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
	g.Coins[player] += amount
}
