package alhambra

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
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanDone(pNum)
}

func (c DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", ErrCouldNotFindPlayer
	}

	return "", g.Done(pNum)
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to end your turn and put all remaining placeable tiles in your reserve"
}

func (g *Game) CanDone(player int) bool {
	return g.CurrentPlayer == player &&
		(g.Phase == PhasePlace || g.Phase == PhaseFinalPlace)
}

func (g *Game) Done(player int) error {
	if !g.CanDone(player) {
		return errors.New("unable to call done right now")
	}

	g.NextPhase()
	return nil
}
