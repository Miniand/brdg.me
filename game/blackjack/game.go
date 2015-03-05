package blackjack

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players []string
	Log     *log.Log

	Dealer        int
	CurrentPlayer int
	CurrentHand   int
	MinBet        int

	Hands map[int][]*Hand
	Money map[int]int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Blackjack"
}

func (g *Game) Identifier() string {
	return "blackjack"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(g.Players) < 2 {
		return errors.New("Blackjack requires at least 2 players")
	}

	g.Log = log.New()
	g.Players = players
	g.MinBet = 5
	g.DealHands()
	return nil
}

func (g *Game) ActivePlayers() []int {
	active := []int{}
	for p := range g.Players {
	}
	return active
}

func (g *Game) PlayerActive(player int) bool {
	return g.Money[player] > 0 || len(g.Hands[player]) > 0
}

func (g *Game) DealHands() {

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
}

func (g *Game) GameLog() *log.Log {
}
