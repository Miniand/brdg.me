package farkle

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (rc RollCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("roll", 0, input)
}

func (rc RollCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.Player] && g.TakenThisRoll && !g.IsFinished()
}

func (rc RollCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	if player != g.Players[g.Player] {
		return "", errors.New("It's not your turn")
	}
	if !g.TakenThisRoll {
		return "", errors.New("You haven't taken any dice this roll")
	}
	if g.IsFinished() {
		return "", errors.New("The game is already finished")
	}
	if len(g.RemainingDice) > 0 {
		g.Roll(len(g.RemainingDice))
	} else {
		g.Roll(6)
	}
	g.TakenThisRoll = false
	return "", nil
}

func (rc RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll{{_b}} to roll remaining dice"
}
