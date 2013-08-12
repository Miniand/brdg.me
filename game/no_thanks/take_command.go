package no_thanks

import (
	"github.com/beefsack/brdg.me/command"
)

type TakeCommand struct{}

func (tc TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 0, input)
}

func (tc TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CurrentlyMoving == player && !g.IsFinished()
}

func (tc TakeCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	return g.Take(player)
}

func (tc TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take{{_b}} to take the card and any chips on it"
}
