package king_of_tokyo

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
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanDone(pNum)
}

func (c DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Done(pNum)
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to finish your turn"
}

func (g *Game) CanDone(player int) bool {
	if g.CurrentPlayer != player {
		return false
	}
	return g.Phase == PhaseBuy
}

func (g *Game) Done(player int) error {
	if !g.CanDone(player) {
		return errors.New("you can't call done at the moment")
	}
	g.NextPhase()
	return nil
}
