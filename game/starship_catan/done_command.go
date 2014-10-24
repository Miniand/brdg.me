package starship_catan

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (c DoneCommand) Parse(input string) []string {
	return command.ParseNamedCommand("done", input)
}

func (c DoneCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanDone(p)
}

func (c DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Done(p)
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to end your turn"
}

func (g *Game) CanDone(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseTradeAndBuild
}

func (g *Game) Done(player int) error {
	if !g.CanDone(player) {
		return errors.New("cannot finish your turn at the moment")
	}
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 2
	g.NewTurn()
	return nil
}
