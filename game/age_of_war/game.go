package age_of_war

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players       []string
	CurrentPlayer int
	Log           *log.Log

	Conquered    map[int]bool
	CastleOwners map[int]int

	CurrentlyAttacking int
	CompletedLines     map[int]bool
	CurrentRoll        []int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		AttackCommand{},
		LineCommand{},
		RollCommand{},
	}
}

func (g *Game) Name() string {
	return "Age of War"
}

func (g *Game) Identifier() string {
	return "age_of_war"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if l := len(players); l < 2 || l > 6 {
		return errors.New("only for 2 to 6 players")
	}
	g.Players = players
	g.Log = log.New()

	g.Conquered = map[int]bool{}
	g.CastleOwners = map[int]int{}

	g.StartTurn()

	return nil
}

func (g *Game) StartTurn() {
	g.CurrentlyAttacking = -1
	g.Roll(7)
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
	for p, name := range g.Players {
		if name == player {
			return p, true
		}
	}
	return 0, false
}
