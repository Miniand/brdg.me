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

var OrthoDirs = []int{
	DirUp,
	DirRight,
	DirDown,
	DirLeft,
}

var OrthoDirNames = map[int]string{
	DirUp:    "up",
	DirRight: "right",
	DirDown:  "down",
	DirLeft:  "left",
}

var DiagDirs = []int{
	DirUp | DirRight,
	DirDown | DirRight,
	DirDown | DirLeft,
	DirUp | DirLeft,
}

var Dirs = append(append([]int{}, OrthoDirs...), DiagDirs...)

type Tiler interface {
	TileAt(loc Loc) (Tile, bool)
}

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

	Board Board

	PlayedPieces map[int]map[int]bool

	CurrentPlayer int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
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

	g.Board = Board{}
	for _, l := range AllLocs {
		g.Board[l] = EmptyTile
	}
	g.PlayedPieces = map[int]map[int]bool{}
	for p := range players {
		g.PlayedPieces[p] = map[int]bool{}
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

func (g *Game) PlayerNum(player string) (int, bool) {
	for pNum, pName := range g.Players {
		if pName == player {
			return pNum, true
		}
	}
	return 0, false
}

func Neighbour(src Tiler, loc Loc, dir int) (Tile, bool) {
	return src.TileAt(loc.Neighbour(dir))
}

func OpenSides(src Tiler, loc Loc) (open map[int]bool) {
	t, ok := src.TileAt(loc)
	if !ok {
		return
	}
	open = map[int]bool{}
	for _, d := range Dirs {
		if nt, ok := Neighbour(src, loc, d); ok && t.Player == nt.Player &&
			t.Type == nt.Type {
			open[d] = true
		}
	}
	return
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = Opponent(g.CurrentPlayer)
}

func Opponent(pNum int) int {
	return (pNum + 1) % 2
}
