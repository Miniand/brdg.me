package sushi_go

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const Dummy = 2

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players    []string
	AllPlayers []string // Includes dummy if needed

	Round int

	Deck    []int
	Hands   map[int][]int
	Playing map[int][]int

	Played map[int][]int
	Points map[int]int

	Controller int // For 2 players, who is controlling the dummy this turn

	Log *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Sushi Go"
}

func (g *Game) Identifier() string {
	return "sushi_go"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	drawCount, ok := PlayerDrawCounts[len(players)]
	if !ok {
		return errors.New("requires between 2 and 5 players")
	}

	g.Log = log.New()
	g.Players = players
	g.AllPlayers = append([]string{}, players...)
	if len(players) == 2 {
		g.AllPlayers = append(g.AllPlayers, "Dummy")
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Because there are only two players, you will be joined by %s",
			g.RenderName(Dummy),
		)))
	}

	g.Deck = Shuffle(Deck())
	g.Hands = map[int][]int{}
	g.Playing = map[int][]int{}
	g.Played = map[int][]int{}
	g.Points = map[int]int{}

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Dealing {{b}}%d{{_b}} cards to each player",
		drawCount,
	)))
	for p := range g.AllPlayers {
		g.Hands[p], g.Deck = g.Deck[:drawCount], g.Deck[drawCount:]
	}
	return nil
}

func (g *Game) StartHand() {
	for p := range g.AllPlayers {
		g.Playing[p] = nil
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
	return []string{}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.AllPlayers[player])
}
