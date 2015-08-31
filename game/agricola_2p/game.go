package agricola_2p

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	Wood = iota
	Stone
	Reed
	Border
	Worker

	Sheep
	Pig
	Cow
	Horse
)

type Game struct {
	Players     []string
	PBoards     [2]*PBoard
	StartPlayer int
	Log         *log.Log
}

func (g *Game) Commands(player string) []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Agricola: All Creatures Big and Small"
}

func (g *Game) Identifier() string {
	return "agricola_2p"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("only for 2 players")
	}
	g.Players = players
	g.Log = log.New()
	g.PBoards = [2]*PBoard{}
	for p := range players {
		g.PBoards[p] = NewPBoard()
	}
	return nil
}

func (g *Game) PlayerList() []string {
	return nil
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

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}
