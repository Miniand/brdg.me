package age_of_war

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players     []string
	CurrentTurn int
	Log         *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
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
