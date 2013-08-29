package liars_dice

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
	"math/rand"
	"time"
)

const (
	START_DICE_COUNT = 5
)

type Game struct {
	Players       []string
	CurrentPlayer int
	PlayerDice    [][]int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Liar's Dice"
}

func (g *Game) Identifier() string {
	return "liars_dice"
}

func (g *Game) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (g *Game) Decode([]byte) error {
	return nil
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) Start(players []string) error {
	// Set players
	if len(players) < 2 || len(players) > 6 {
		return errors.New("Liar's Dice must be between 2 and 6 players")
	}
	g.Players = players
	// Set a random first player
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.CurrentPlayer = r.Int() % len(g.Players)
	// Initialise dice
	g.PlayerDice = make([][]int, len(g.Players))
	for pNum, _ := range g.Players {
		g.PlayerDice[pNum] = make([]int, START_DICE_COUNT)
	}
	// Kick off the first round
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.RollDice()
}

func (g *Game) RollDice() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for pNum, _ := range g.Players {
		for d, _ := range g.PlayerDice[pNum] {
			g.PlayerDice[pNum][d] = (r.Int() % 6) + 1
		}
	}
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
