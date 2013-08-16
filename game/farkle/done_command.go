package farkle

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (dc DoneCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("done", 0, input)
}

func (dc DoneCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.Player] && g.TurnScore > 0 && !g.IsFinished()
}

func (dc DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	if player != g.Players[g.Player] {
		return "", errors.New("It's not your turn")
	}
	if g.TurnScore == 0 {
		return "", errors.New("You haven't scored anything yet")
	}
	if g.IsFinished() {
		return "", errors.New("The game is already finished")
	}
	g.NextPlayer()
	return "", nil
}

func (dc DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to take the points and finish your turn"
}
