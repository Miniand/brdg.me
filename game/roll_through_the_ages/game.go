package roll_through_the_ages

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	PhaseRoll = iota
	PhaseExtraRoll
	PhaseCollect
	PhaseResolve
	PhaseBuild
	PhaseBuy
	PhaseDiscard
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players          []string
	CurrentPlayer    int
	Phase            int
	Boards           []*PlayerBoard
	RolledDice       []int
	KeptDice         []int
	RemainingRolls   int
	RemainingCoins   int
	RemainingWorkers int
	Log              *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		RollCommand{},
		KeepCommand{},
		TakeCommand{},
	}
}

func (g *Game) Name() string {
	return "Roll Through the Ages"
}

func (g *Game) Identifier() string {
	return "roll_through_the_ages"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	l := len(players)
	if l < 1 || l > 4 {
		return errors.New("Roll Through the Ages is 1-4 player")
	}
	g.Players = players
	g.Boards = make([]*PlayerBoard, l)
	for i := 0; i < l; i++ {
		g.Boards[i] = NewPlayerBoard()
	}
	g.Log = log.New()
	g.StartTurn()
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
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) StartTurn() {
	g.Phase = PhaseRoll
	g.NewRoll(g.Boards[g.CurrentPlayer].Cities())
	g.RemainingRolls = 2
}

func (g *Game) RollExtraPhase() {
	g.Phase = PhaseExtraRoll
	if !g.Boards[g.CurrentPlayer].Developments[DevelopmentLeadership] {
		g.CollectPhase()
	}
}

func (g *Game) CollectPhase() {
	g.Phase = PhaseCollect
	g.KeptDice = append(g.RolledDice, g.KeptDice...)
	g.RolledDice = []int{}
	// Collect goods and food
	cp := g.CurrentPlayer
	hasFoodOrWorkersDice := false
	goods := 0
	for _, d := range g.KeptDice {
		switch d {
		case DiceFood:
			g.Boards[cp].Food += 3 + g.Boards[cp].FoodModifier()
		case DiceGood:
			goods += 1
		case DiceSkull:
			goods += 2
		case DiceWorkers:
			g.RemainingWorkers += 3 + g.Boards[cp].WorkerModifier()
		case DiceFoodOrWorkers:
			hasFoodOrWorkersDice = true
		case DiceCoins:
			g.RemainingCoins += g.Boards[cp].CoinsDieValue()
		}
	}
	if !hasFoodOrWorkersDice {
		g.PhaseResolve()
	}
}

func (g *Game) PhaseResolve() {
	g.Phase = PhaseResolve
	// Feed cities
	// Resolve disasters
	g.BuildPhase()
}

func (g *Game) BuildPhase() {
	g.Phase = PhaseBuild
}

func (g *Game) BuyPhase() {
	g.Phase = PhaseBuy
}

func (g *Game) DiscardPhase() {
	g.Phase = PhaseDiscard
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	g.StartTurn()
}

func (g *Game) PlayerNum(player string) (int, error) {
	for pNum, p := range g.Players {
		if player == p {
			return pNum, nil
		}
	}
	return 0, fmt.Errorf("could not find a player by the name %s", player)
}

func ContainsInt(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}
