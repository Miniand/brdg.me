package no_thanks

import "github.com/Miniand/brdg.me/command"

type TakeCommand struct{}

func (tc TakeCommand) Name() string { return "take" }

func (tc TakeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	return "", g.Take(player)
}

func (tc TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take{{_b}} to take the card and any chips on it"
}

func (g *Game) CanTake(player string) bool {
	return g.CurrentlyMoving == player && !g.IsFinished()
}
