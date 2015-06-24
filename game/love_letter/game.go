package love_letter

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players                []string
	Log                    *log.Log
	Deck, Removed, Discard []int
	Hands                  [][]int
	Points                 []int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Love Letter"
}

func (g *Game) Identifier() string {
	return "love_letter"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	l := len(players)
	if l < 2 || l > 4 {
		return errors.New("only for 2 to 4 players")
	}
	g.Players = players
	g.Log = log.New()
	g.Points = make([]int, l)
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	deck := helper.IntShuffle(Deck)
	remove := 1
	l := len(g.Players)
	if l == 2 {
		remove = 4
	}
	g.Deck, g.Removed = deck[remove:], deck[:remove]
	g.Hands = make([][]int, l)
	for p := range g.Players {
		g.Hands[p] = []int{}
		g.DrawCard(p)
	}
}

func (g *Game) DrawCard(player int) {
	var card int
	if len(g.Deck) > 0 {
		card, g.Deck = g.Deck[0], g.Deck[1:]
	} else {
		card = g.Removed[0]
	}
	g.Hands[player] = append(g.Hands[player], card)
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
