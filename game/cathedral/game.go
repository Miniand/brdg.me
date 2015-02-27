package cathedral

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	DirUp = 1 << iota
	DirRight
	DirDown
	DirLeft
)

type Game struct {
	Players []string
	Log     *log.Log

	Board [10][10]Tile

	CurrentPlayer int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Cathedral"
}

func (g *Game) Identifier() string {
	return "cathedral"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Cathedral is two player")
	}
	g.Players = players
	g.Log = log.New()

	board := [10][10]Tile{}
	for i := 0; i < 10; i++ {
		board[i] = [10]Tile{}
		for j := 0; j < 10; j++ {
			board[i][j] = Tile{
				Player: NoPlayer,
			}
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
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}
