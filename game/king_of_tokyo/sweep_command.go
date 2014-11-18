package king_of_tokyo

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type SweepCommand struct{}

func (c SweepCommand) Parse(input string) []string {
	return command.ParseNamedCommand("sweep", input)
}

func (c SweepCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanSweep(pNum)
}

func (c SweepCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Sweep(pNum)
}

func (c SweepCommand) Usage(player string, context interface{}) string {
	return fmt.Sprintf(
		"{{b}}sweep{{_b}} to spend %s to discard all cards for purchase and draw new ones",
		RenderEnergy(2),
	)
}

func (g *Game) CanSweep(player int) bool {
	if g.CurrentPlayer != player {
		return false
	}
	return g.Phase == PhaseBuy && len(g.Buyable) > 0 &&
		g.Boards[player].Energy >= 2
}

func (g *Game) Sweep(player int) error {
	if !g.CanSweep(player) {
		return errors.New("you can't sweep at the moment")
	}
	n := 3
	if l := len(g.Deck); l < n {
		n = l
	}
	g.Boards[player].Energy -= 2
	g.Discard = append(g.Discard, g.Buyable...)
	g.Buyable = g.Deck[:n]
	g.Deck = g.Deck[n:]
	return nil
}
