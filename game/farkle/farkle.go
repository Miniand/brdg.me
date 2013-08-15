package farkle

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Game struct {
	Players       []string
	FirstPlayer   int
	Player        int
	Scores        map[int]int
	TurnScore     int
	RemainingDice []int
	TakenThisRoll bool
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Farkle"
}

func (g *Game) Identifier() string {
	return "farkle"
}

func (g *Game) Encode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(g)
	return buf.Bytes(), err
}

func (g *Game) Decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(g)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	output := bytes.NewBufferString("")
	renderedDice := make([]string, len(g.RemainingDice))
	for i, d := range g.RemainingDice {
		renderedDice[i] = die.Render(d)
	}
	output.WriteString(strings.Join(renderedDice, " "))
	return output.String(), nil
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 {
		return errors.New("Farkle requires at least two players")
	}
	g.Players = players
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Player = r.Int() % len(g.Players)
	g.FirstPlayer = g.Player
	g.StartTurn()
	return nil
}

func (g *Game) StartTurn() {
	g.TurnScore = 0
	g.TakenThisRoll = false
	g.Roll(6)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.Player]}
}

func (g *Game) Roll(n int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.RemainingDice = make([]int, n)
	for i := 0; i < n; i++ {
		g.RemainingDice[i] = r.Int()%6 + 1
	}
	sort.IntSlice(g.RemainingDice).Sort()
}
