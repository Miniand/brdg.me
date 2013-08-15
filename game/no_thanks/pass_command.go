package no_thanks

import (
	"github.com/Miniand/brdg.me/command"
)

type PassCommand struct{}

func (pc PassCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("pass", 0, input)
}

func (pc PassCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CurrentlyMoving == player && g.PlayerChips[player] > 0 &&
		!g.IsFinished()
}

func (pc PassCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	return "", g.Pass(player)
}

func (pc PassCommand) Usage(player string, context interface{}) string {
	return "{{b}}pass{{_b}} to spend a chip to pass on this card"
}
