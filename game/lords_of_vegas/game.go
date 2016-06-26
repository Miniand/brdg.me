package lords_of_vegas

import (
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	CasinoNone = iota
	CasinoAlbion
	CasinoSphynx
	CasinoVega
	CasinoTivoli
	CasinoPioneer
	CasinoTheStrip
)

type Game struct {
	Players  []string
	Log      *log.Log
	Finished bool
}

func (g *Game) Commands(player string) []command.Command {
	commands := []command.Command{}
	return commands
}

func (g *Game) Name() string {
	return "Lords of Vegas"
}

func (g *Game) Identifier() string {
	return "lords_of_vegas"
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
	return g.Finished
}

func (g *Game) Winners() []string {
	return nil
}

func (g *Game) WhoseTurn() []string {
	return nil
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}
