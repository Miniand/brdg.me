package hive

import (
	"encoding/json"
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/grid"
	"github.com/Miniand/brdg.me/game/grid/hex"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players []string
	Log     *log.Log
	Board   hex.Grid
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Must be 2 players")
	}
	g.Players = players
	g.Log = log.New()
	g.Board = hex.Grid{}
	return nil
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	g.Board.SetTile(grid.Loc{0, 0}, &Tile{TILE_BEETLE, 0})
	g.Board.SetTile(grid.Loc{1, 0}, &Tile{TILE_QUEEN_BEE, 1})
	g.Board.SetTile(grid.Loc{1, 1}, &Tile{TILE_GRASSHOPPER, 0})
	g.Board.SetTile(grid.Loc{0, 1}, &Tile{TILE_SPIDER, 1})
	return render.RenderHexGrid(g.Board, 2), nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) WhoseTurn() []string {
	return g.Players
}

func (g *Game) Name() string {
	return "Hive"
}

func (g *Game) Identifier() string {
	return "hive"
}

// Encode to a string
func (g *Game) Encode() ([]byte, error) {
	return json.Marshal(g)
}

// Decode from a string
func (g *Game) Decode(data []byte) error {
	return json.Unmarshal(data, g)
}
