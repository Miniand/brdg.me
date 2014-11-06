package king_of_tokyo

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	PhaseRoll = iota
	PhaseBuy
)

const (
	LocationOutside  = -1
	LocationTokyo    = 0
	LocationTokyoBay = 1
)

type Game struct {
	Players        []string
	CurrentPlayer  int
	Phase          int
	CurrentRoll    []int
	RemainingRolls int
	Buyable        []CardBase
	Deck           []CardBase
	Boards         []*PlayerBoard
	Tokyo          []int
	Log            *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "King of Tokyo"
}

func (g *Game) Identifier() string {
	return "king_of_tokyo"
}

func (g *Game) PlayerLocation(player int) int {
	for loc, p := range g.Tokyo {
		if p == player {
			return loc
		}
	}
	return LocationOutside
}

func (g *Game) PlayersInsideTokyo() []int {
	players := []int{}
	for _, p := range g.Tokyo {
		if p != -1 {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) PlayersOutsideTokyo() []int {
	players := []int{}
	inside := g.PlayersInsideTokyo()
	for p, _ := range g.Players {
		outside := true
		for _, in := range inside {
			if p == in {
				outside = false
				break
			}
		}
		if outside {
			players = append(players, p)
		}
	}
	return players
}

func RegisterGobTypes() {
	for _, c := range Deck() {
		gob.Register(c)
	}
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	playerCount := len(players)
	if playerCount < 2 || playerCount > 6 {
		return errors.New("requires between 2 and 6 players")
	}
	g.Players = players
	g.Log = log.New()
	deck := Shuffle(Deck())
	g.Buyable = deck[:3]
	g.Deck = deck[3:]
	locations := 1
	if playerCount > 4 {
		locations = 2
	}
	g.Tokyo = make([]int, locations)
	for i, _ := range g.Tokyo {
		g.Tokyo[i] = -1
	}
	g.Boards = make([]*PlayerBoard, playerCount)
	for p, _ := range g.Players {
		g.Boards[p] = NewPlayerBoard()
	}
	g.StartTurn()
	return nil
}

func (g *Game) StartTurn() {
	g.CurrentRoll = RollDice(6)
	g.LogRoll(g.CurrentRoll, []int{})
	g.RemainingRolls = 2
}

func (g *Game) LogRoll(rolled, kept []int) {
	diceStr := []string{}
	for _, d := range rolled {
		diceStr = append(diceStr, render.Bold(RenderDie(d)))
	}
	for _, d := range kept {
		diceStr = append(diceStr, RenderDie(d))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s rolled  %s",
		g.RenderName(g.CurrentPlayer),
		strings.Join(diceStr, "  "),
	)))
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

func (g *Game) EliminatedPlayerList() []string {
	return []string{}
}
