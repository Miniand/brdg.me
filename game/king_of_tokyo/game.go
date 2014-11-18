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
	LocationOutside   = -1
	LocationTokyoCity = 0
	LocationTokyoBay  = 1
)

type Game struct {
	Players        []string
	CurrentPlayer  int
	Phase          int
	CurrentRoll    []int
	RemainingRolls int
	Buyable        []CardBase
	Deck           []CardBase
	Discard        []CardBase
	Boards         []*PlayerBoard
	Tokyo          []int
	Log            *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		RollCommand{},
		KeepCommand{},
		BuyCommand{},
		SweepCommand{},
		DoneCommand{},
	}
}

func (g *Game) NextPhase() {
	switch g.Phase {
	case PhaseRoll:
		g.BuyPhase()
	case PhaseBuy:
		g.NextTurn()
	}
}

func (g *Game) RollPhase() {
	g.Phase = PhaseRoll
	g.CurrentRoll = RollDice(6)
	g.LogRoll(g.CurrentRoll, []int{})
	g.RemainingRolls = 2
}

func (g *Game) BuyPhase() {
	// Handle dice
	diceCounts := map[int]int{}
	for _, d := range g.CurrentRoll {
		diceCounts[d] += 1
	}
	for _, d := range Dice {
		count := diceCounts[d]
		if count == 0 {
			continue
		}
		switch d {
		case Die1, Die2, Die3:
			if count >= 3 {
				g.Boards[g.CurrentPlayer].VP += d + count - 2
			}
		case DieEnergy:
			g.Boards[g.CurrentPlayer].Energy += count
		case DieAttack:

		case DieHeal:
			g.Boards[g.CurrentPlayer].Health += count
		}
	}
	// Start buy phase
	g.Phase = PhaseBuy
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
	g.Discard = []CardBase{}
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
	g.RollPhase()
	return nil
}

func (g *Game) NextTurn() {
	if !g.IsFinished() {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		for g.Boards[g.CurrentPlayer].Health <= 0 {
			g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		}
		g.RollPhase()
	}
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

func (g *Game) CheckGameEnd() (isFinished bool, winners []string) {
	winners = []string{}
	remainingPlayers := []string{}
	for p, b := range g.Boards {
		if b.Health > 0 {
			if b.VP >= 20 {
				return true, []string{g.Players[p]}
			}
			remainingPlayers = append(remainingPlayers, g.Players[p])
		}
	}
	if len(remainingPlayers) < 2 {
		isFinished = true
		winners = remainingPlayers
	}
	return
}

func (g *Game) IsFinished() bool {
	isFinished, _ := g.CheckGameEnd()
	return isFinished
}

func (g *Game) Winners() []string {
	_, winners := g.CheckGameEnd()
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) EliminatedPlayerList() []string {
	eliminated := []string{}
	for p, b := range g.Boards {
		if b.Health <= 0 {
			eliminated = append(eliminated, g.Players[p])
		}
	}
	return eliminated
}

func (g *Game) PlayerNum(player string) (int, error) {
	return helper.StringInStrings(player, g.Players)
}

func ContainsInt(i int, s []int) bool {
	for _, si := range s {
		if si == i {
			return true
		}
	}
	return false
}
