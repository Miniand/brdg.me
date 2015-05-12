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
	Hands   [][]int
	Playing map[int][]int

	Played map[int][]int
	Points map[int]int

	Controller int // For 2 players, who is controlling the dummy this turn

	Log *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
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
	_, ok := PlayerDrawCounts[len(players)]
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
	g.Playing = map[int][]int{}
	g.Played = map[int][]int{}
	g.Points = map[int]int{}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.Round++
	g.Hands = make([][]int, len(g.AllPlayers))
	drawCount := PlayerDrawCounts[len(g.Players)]
	passDir := "left"
	if g.Round == 2 {
		passDir = "right"
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Starting round {{b}}%d{{_b}}, hands will be passed to the {{b}}%s{{_b}}.  Dealing {{b}}%d{{_b}} cards to each player",
		g.Round,
		passDir,
		drawCount,
	)))
	for p := range g.AllPlayers {
		g.Hands[p], g.Deck = g.Deck[:drawCount], g.Deck[drawCount:]
	}
	g.StartHand()
}

func (g *Game) StartHand() {
}

func (g *Game) EndHand() {
	// Play cards
	for p := range g.AllPlayers {
		g.Played[p] = append(g.Played[p], g.Playing[p]...)
		if len(g.Playing) == 2 {
			// Use chopsticks.
			if i, ok := Contains(CardChopsticks, g.Played[p]); ok {
				g.Hands[p] = append(g.Hands[p], CardChopsticks)
				g.Played[p] = append(g.Played[p][:i], g.Played[p][i+1:]...)
			}
		}
		g.Playing[p] = nil
	}
	// End round if we're out of cards
	if len(g.Hands[0]) == 0 {
		g.EndRound()
		return
	}
	// Pass hands
	if len(g.Players) == 2 {
		g.Log.Add(log.NewPublicMessage("Players are swapping hands"))
		g.Hands[0], g.Hands[1] = g.Hands[1], g.Hands[0]
		// Next player controls the dummy
		g.Controller = (g.Controller + 1) % 2
	} else if g.Round%2 == 1 {
		g.Log.Add(log.NewPublicMessage("Passing hands to the {{b}}left{{_b}}"))
		extra := g.Hands[0]
		g.Hands = append(g.Hands[1:], extra)
	} else {
		g.Log.Add(log.NewPublicMessage("Passing hands to the {{b}}right{{_b}}"))
		l := len(g.Hands)
		extra := g.Hands[l-1]
		g.Hands = append([][]int{extra}, g.Hands[:l-1]...)
	}
}

func (g *Game) EndRound() {
	if g.Round < 3 {
		g.StartRound()
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Round == 3 && len(g.Hands[0]) == 0 && g.Playing[0] == nil
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	whose := []string{}
	commands := g.Commands()
	for _, pName := range g.Players {
		if len(command.AvailableCommands(pName, g, commands)) > 0 {
			whose = append(whose, pName)
		}
	}
	return whose
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

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.AllPlayers[player])
}
