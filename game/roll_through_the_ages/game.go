package roll_through_the_ages

import (
	"bytes"
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
	Round            int
	FinalRound       int
	Log              *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		RollCommand{},
		TakeCommand{},
		BuildCommand{},
		BuyCommand{},
		DiscardCommand{},
		NextCommand{},
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
	g.Round = 1
	if len(players) == 1 {
		g.FinalRound = 10
	}
	g.Log = log.New()
	g.StartTurn()
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.FinalRound != 0 && g.Round > g.FinalRound
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []int{}
	winningScore := 0
	for p, _ := range g.Players {
		score := g.Boards[p].Score()
		if score > winningScore {
			winners = []int{}
			winningScore = score
		}
		if score == winningScore {
			winners = append(winners, p)
		}
	}
	if len(winners) < 2 {
		return g.PlayerNumsToNames(winners)
	}
	// There's a tie, goods value is tie breaker
	goodsWinners := []int{}
	goodsScore := 0
	for p, _ := range winners {
		score := g.Boards[p].GoodsValue()
		if score > goodsScore {
			goodsWinners = []int{}
			goodsScore = score
		}
		if score == goodsScore {
			goodsWinners = append(goodsWinners, p)
		}
	}
	return g.PlayerNumsToNames(goodsWinners)
}

func (g *Game) PlayerNumsToNames(players []int) []string {
	names := make([]string, len(players))
	for i, p := range players {
		names[i] = g.Players[p]
	}
	return names
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
	g.Boards[cp].GainGoods(goods)
	if !hasFoodOrWorkersDice {
		g.PhaseResolve()
	}
}

func (g *Game) PhaseResolve() {
	g.Phase = PhaseResolve
	cp := g.CurrentPlayer
	// Feed cities
	if cities := g.Boards[cp].Cities(); g.Boards[cp].Food >= cities {
		g.Boards[cp].Food -= cities
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s fed {{b}}%d{{_b}} cities`,
			g.RenderName(cp),
			cities,
		)))
	} else {
		// Famine
		famine := cities - g.Boards[cp].Food
		g.Boards[cp].Food = 0
		g.Boards[cp].Disasters += famine
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`Famine! %s takes {{b}}%d disaster points{{_b}}`,
			g.RenderName(cp),
			famine,
		)))
	}
	// Resolve disasters
	skulls := 0
	for _, d := range g.KeptDice {
		if d == DiceSkull {
			skulls += 1
		}
	}
	switch skulls {
	case 0, 1:
		break
	case 2:
		if g.Boards[cp].Developments[DevelopmentIrrigation] {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`%s avoids a drought with their irrigation development`,
				g.RenderName(cp),
			)))
		} else {
			g.Boards[cp].Disasters += 2
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`Drought! %s takes {{b}}2 disaster points{{_b}}`,
				g.RenderName(cp),
			)))
		}
	case 3:
		buf := bytes.NewBufferString("Pestilence!")
		for p, _ := range g.Players {
			if p == cp {
				continue
			}
			if g.Boards[p].Developments[DevelopmentMedicine] {
				buf.WriteString(fmt.Sprintf(
					"\n  %s avoids pestilence with their medicine development",
					g.RenderName(p),
				))
			} else {
				g.Boards[cp].Disasters += 3
				buf.WriteString(fmt.Sprintf(
					"\n  %s takes {{b}}3 disaster points{{_b}}",
					g.RenderName(p),
				))
			}
		}
		g.Log.Add(log.NewPublicMessage(buf.String()))
	case 4:
		if g.Boards[cp].HasBuilt(MonumentGreatWall) {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`%s avoids an invasion with their wall`,
				g.RenderName(cp),
			)))
		} else {
			g.Boards[cp].Disasters += 4
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`Invasion! %s takes {{b}}4 disaster points{{_b}}`,
				g.RenderName(cp),
			)))
		}
	default:
		for _, good := range Goods {
			g.Boards[cp].Goods[good] = 0
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`Revolt! %s loses {{b}}all of their goods{{_b}}`,
			g.RenderName(cp),
		)))
	}
	g.BuildPhase()
}

func (g *Game) BuildPhase() {
	g.Phase = PhaseBuild
	if g.RemainingWorkers == 0 {
		g.BuyPhase()
	}
}

func (g *Game) BuyPhase() {
	g.Phase = PhaseBuy
	b := g.Boards[g.CurrentPlayer]
	if g.RemainingCoins == 0 && b.GoodsNum() == 0 &&
		(!b.Developments[DevelopmentGranaries] || b.Food == 0) {
		g.DiscardPhase()
	}
}

func (g *Game) DiscardPhase() {
	g.Phase = PhaseDiscard
	if g.Boards[g.CurrentPlayer].GoodsNum() <= 6 ||
		g.Boards[g.CurrentPlayer].Developments[DevelopmentCaravans] {
		g.NextTurn()
	}
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
