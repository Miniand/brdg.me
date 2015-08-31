package zombie_dice

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

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players []string
	Log     *log.Log

	CurrentTurn    int
	Scores         []int
	Cup            []Dice
	RollOffPlayers map[int]bool
	Finished       bool
	CurrentRoll    DiceResultList
	Kept           DiceResultList
	RoundBrains    int
	RoundShotguns  int
}

func (g *Game) Commands(player string) []command.Command {
	commands := []command.Command{}
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return commands
	}
	if g.CanRoll(pNum) {
		commands = append(commands, RollCommand{})
	}
	if g.CanKeep(pNum) {
		commands = append(commands, KeepCommand{})
	}
	return commands
}

func (g *Game) Name() string {
	return "Zombie Dice"
}

func (g *Game) Identifier() string {
	return "zombie_dice"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	pLen := len(players)
	if pLen < 2 {
		return errors.New("requires at least 2 players")
	}

	g.Players = players
	g.Log = log.New()

	g.Scores = make([]int, pLen)
	g.StartTurn()

	return nil
}

func (g *Game) ShakeCup() {
	l := len(g.Cup)
	shaken := make([]Dice, l)
	for i, p := range rnd.Perm(l) {
		shaken[i] = g.Cup[p]
	}
	g.Cup = shaken
}

func (g *Game) TakeDice(n int) []Dice {
	if n < 0 {
		panic("Must have more than 0")
	}
	dice := []Dice{}
	if n == 0 {
		return dice
	}

	if len(g.Cup) < n {
		g.Log.Add(log.NewPublicMessage(
			"Not enough dice remaining, returning kept dice to the cup"))
		for _, d := range g.Kept {
			g.Cup = append(g.Cup, d.Dice)
		}
		g.Kept = DiceResultList{}
		g.ShakeCup()
	}

	dice, g.Cup = g.Cup[:n], g.Cup[n:]
	return dice
}

func (g *Game) StartTurn() {
	g.Cup = AllDice()
	g.ShakeCup()
	g.Kept = DiceResultList{}
	g.CurrentRoll = DiceResultList{}
	g.RoundBrains = 0
	g.RoundShotguns = 0
	g.Roll()
}

func (g *Game) NextPlayer() {
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)
	if g.CurrentTurn == 0 {
		// Check for game end
		score, leaders := g.Leaders()
		if score >= 13 {
			if len(leaders) == 1 {
				g.Finished = true
				return
			}
			// Roll off!
			g.RollOffPlayers = map[int]bool{}
			parts := []string{}
			for _, l := range leaders {
				g.RollOffPlayers[l] = true
				parts = append(parts, g.PlayerName(l))
			}
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"It's a tied score of {{b}}%d{{_b}} between %s, tie breaker round!",
				score,
				render.CommaList(parts),
			)))
		}
	}
	if g.RollOffPlayers != nil && !g.RollOffPlayers[g.CurrentTurn] {
		g.NextPlayer()
	} else {
		g.StartTurn()
	}
}

func (g *Game) Roll() {
	dice := g.CurrentRoll.Dice()
	diceLen := len(dice)
	if diceLen < 3 {
		dice = append(dice, g.TakeDice(3-diceLen)...)
	}
	drl := RollDice(dice)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s rolled %s",
		g.PlayerName(g.CurrentTurn),
		drl,
	)))

	run := DiceResultList{}
	newBrains := 0
	wasShot := false
	for _, dr := range drl {
		switch dr.Face {
		case Brain:
			newBrains++
			g.Kept = append(g.Kept, dr)
		case Shotgun:
			g.RoundShotguns++
			g.Kept = append(g.Kept, dr)
			wasShot = true
		case Footprints:
			run = append(run, dr)
		}
	}
	if g.RoundShotguns >= 3 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s got shot three times and lost {{b}}%d{{_b}} brains!",
			g.PlayerName(g.CurrentTurn),
			g.RoundBrains,
		)))
		g.NextPlayer()
		return
	} else if wasShot {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s has {{b}}%d{{_b}} health remaining",
			g.PlayerName(g.CurrentTurn),
			3-g.RoundShotguns,
		)))
	}
	g.RoundBrains += newBrains

	g.CurrentRoll = run
}

func (g *Game) Keep() {
	g.Scores[g.CurrentTurn] += g.RoundBrains
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s kept {{b}}%d{{_b}} brains, now has {{b}}%d{{_b}}!",
		g.PlayerName(g.CurrentTurn),
		g.RoundBrains,
		g.Scores[g.CurrentTurn],
	)))
	g.NextPlayer()
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Finished
}

func (g *Game) Leaders() (score int, players []int) {
	players = []int{0}
	for p, _ := range g.Players {
		if g.Scores[p] > score {
			score = g.Scores[p]
			players = []int{}
		}
		if g.Scores[p] == score {
			players = append(players, p)
		}
	}
	return score, players
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	_, leaders := g.Leaders()
	return []string{g.Players[leaders[0]]}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	return []string{g.Players[g.CurrentTurn]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for p, name := range g.Players {
		if name == player {
			return p, true
		}
	}
	return 0, false
}
