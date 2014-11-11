package bang_dice

import (
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
	Players []string
	Log     *log.Log
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
	g.Players = players
	g.Log = log.New()
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
	return []string{}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) EliminatedPlayerList() []string {
	return []string{}
}
