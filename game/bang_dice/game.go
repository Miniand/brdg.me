package bang_dice

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
	RoleSheriff = iota
	RoleDeputy
	RoleOutlaw
	RoleRenegade
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

var Roles = []int{
	RoleSheriff,
	RoleRenegade,
	RoleOutlaw,
	RoleOutlaw,
	RoleDeputy,
	RoleOutlaw,
	RoleDeputy,
	RoleRenegade,
}

type Game struct {
	Players            []string
	Roles, Chars, Life []int
	Log                *log.Log
	CurrentTurn        int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Bang! The Dice Game"
}

func (g *Game) Identifier() string {
	return "bang_dice"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	l := len(players)
	if l < 4 || l > 8 {
		return errors.New("only for 4 to 8 players")
	}
	g.Players = players
	g.Log = log.New()

	g.Chars = make([]int, l)
	for i, c := range rnd.Perm(len(Chars))[:l] {
		g.Chars[i] = c
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s is {{b}}%s{{_b}} with {{b}}%d{{_b}} life.  %s",
			g.PlayerName(i),
			Chars[c].Name(),
			Chars[c].StartingLife(),
			Chars[c].Description(),
		)))
	}

	g.Roles = make([]int, l)
	for i, p := range rnd.Perm(l) {
		g.Roles[i] = Roles[p]
		if Roles[p] == RoleSheriff {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"{{b}}%s is the Sherrif{{_b}}, they start with an extra life and take the first turn.",
				g.PlayerName(i),
			)))
			g.CurrentTurn = i
		}
	}

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
	return []string{g.Players[g.CurrentTurn]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) EliminatedPlayerList() []string {
	return []string{}
}
