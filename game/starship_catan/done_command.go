package starship_catan

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (c DoneCommand) Name() string { return "done" }

func (c DoneCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
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
