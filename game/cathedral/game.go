package cathedral

import (
	"errors"
	"math/rand"
	"time"

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

var OrthoDirs = []int{
	DirUp,
	DirRight,
	DirDown,
	DirLeft,
}

var DiagDirs = []int{
	DirUp | DirRight,
	DirDown | DirRight,
	DirDown | DirLeft,
	DirUp | DirLeft,
}

var Dirs = append(append([]int{}, OrthoDirs...), DiagDirs...)

func DirInv(dir int) int {
	var inv int
	if dir&DirUp > 0 {
		inv = inv | DirDown
	}
	if dir&DirRight > 0 {
		inv = inv | DirLeft
	}
	if dir&DirDown > 0 {
		inv = inv | DirUp
	}
	if dir&DirLeft > 0 {
		inv = inv | DirRight
	}
	return inv
}

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

	g.Board = [10][10]Tile{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		g.Board[i] = [10]Tile{}
		for j := 0; j < 10; j++ {
			g.Board[i][j] = Tile{
				Player: (r.Int() % 4) - 1,
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

func (g *Game) TileAt(x, y int) (Tile, bool) {
	if x < 0 || x >= 10 || y < 0 || y >= 10 {
		return Tile{}, false
	}
	return g.Board[y][x], true
}

func (g *Game) Neighbour(x, y, dir int) (Tile, bool) {
	if dir&DirUp == DirUp {
		y--
	}
	if dir&DirRight == DirRight {
		x++
	}
	if dir&DirDown == DirDown {
		y++
	}
	if dir&DirLeft == DirLeft {
		x--
	}
	return g.TileAt(x, y)
}

func (g *Game) OpenSides(x, y int) (open map[int]bool) {
	t, ok := g.TileAt(x, y)
	if !ok {
		return
	}
	open = map[int]bool{}
	for _, d := range Dirs {
		if nt, ok := g.Neighbour(x, y, d); ok && t.Player == nt.Player &&
			t.Type == nt.Type {
			open[d] = true
		}
	}
	return
}
