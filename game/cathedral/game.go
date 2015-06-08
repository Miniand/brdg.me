package cathedral

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
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

	NoOpenTiles bool
	Finished    bool
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
	return g.Finished
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	p0 := g.RemainingPieceSize(0)
	p1 := g.RemainingPieceSize(1)
	if p0 < p1 {
		return []string{g.Players[0]}
	} else if p1 < p0 {
		return []string{g.Players[1]}
	}
	return g.Players
}

func (g *Game) WhoseTurn() []string {
	if g.NoOpenTiles {
		players := []string{}
		for p, pName := range g.Players {
			if g.CanPlaySomething(p, LocFilterPlayable) {
				players = append(players, pName)
			}
		}
		return players
	}
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
	opponent := Opponent(g.CurrentPlayer)
	if g.CanPlaySomething(opponent, LocFilterPlayable) {
		g.CurrentPlayer = opponent
	} else if !g.CanPlaySomething(g.CurrentPlayer, LocFilterPlayable) {
		g.Finished = true
		buf := bytes.NewBufferString(render.Bold(
			"The game is finished, remaining piece size is as follows:",
		))
		for p := range g.Players {
			buf.WriteString(fmt.Sprintf(
				"\n%s - {{b}}%d{{_b}}",
				g.PlayerName(p),
				g.RemainingPieceSize(p),
			))
		}
		g.Log.Add(log.NewPublicMessage(buf.String()))
	}
}

func (g *Game) RemainingPieceSize(player int) int {
	sum := 0
	for pNum, piece := range Pieces[player] {
		if !g.PlayedPieces[player][pNum] {
			sum += len(piece.Positions)
		}
	}
	return sum
}

type LocFilter func(g *Game, player int, loc Loc) bool

func LocFilterPlayable(g *Game, player int, loc Loc) bool {
	t := g.Board[loc]
	return t.Player == NoPlayer && (t.Owner == NoPlayer || t.Owner == player)
}

func LocFilterOpen(g *Game, player int, loc Loc) bool {
	t := g.Board[loc]
	return t.Player == NoPlayer && t.Owner == NoPlayer
}

func (g *Game) CanPlaySomething(player int, filter LocFilter) bool {
	for _, l := range AllLocs {
		if !filter(g, player, l) {
			continue
		}
		// Try to play the easiest one first
		for i := len(Pieces[player]) - 1; i >= 0; i-- {
			if g.PlayedPieces[player][i] {
				continue
			}
			dirs := OrthoDirs
			if !Pieces[player][i].Directional {
				dirs = []int{DirDown}
			}
			for _, dir := range dirs {
				if ok, _ := g.CanPlayPiece(player, i, l, dir); ok {
					return true
				}
			}
		}
	}
	return false
}

func Opponent(pNum int) int {
	return (pNum + 1) % 2
}
