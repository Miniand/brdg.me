package red7

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

	Finished bool

	Deck        []int
	Hands       [][]int
	Palettes    [][]int
	ScoredCards [][]int
	Eliminated  []bool
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
}

func (g *Game) Name() string {
	return "Red7"
}

func (g *Game) Identifier() string {
	return "red7"
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
		return errors.New("only 2-4 players")
	}
	g.Players = players
	g.Log = log.New()

	g.Hands = make([][]int, l)
	g.Palettes = make([][]int, l)
	g.ScoredCards = make([][]int, l)

	g.StartRound()

	return nil
}

func (g *Game) StartRound() {
	l := len(g.Players)

	// Add hands and palettes back to the deck.
	for p := range g.Players {
		g.Deck = append(g.Deck, g.Hands[p]...)
		g.Deck = append(g.Deck, g.Palettes[p]...)
	}
	g.Hands = make([][]int, l)
	g.Palettes = make([][]int, l)

	if len(g.Deck) < l*8 {
		// End of the game, not enough cards to deal new hand.
		g.Finished = true
	}

	g.Deck = helper.IntShuffle(Deck)

	// Deal hands and new palettes.
	for p := range g.Players {
		g.Hands[p] = helper.IntSort(g.Deck[:7])
		g.Palettes[p] = g.Deck[7:8]
		g.Deck = g.Deck[8:]
	}
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
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for p, pName := range g.Players {
		if pName == player {
			return p, true
		}
	}
	return 0, false
}
