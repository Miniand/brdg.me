package no_thanks

import "github.com/Miniand/brdg.me/command"

type PassCommand struct{}

func (pc PassCommand) Name() string { return "pass" }

func (pc PassCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	return "", g.Pass(player)
}

func (pc PassCommand) Usage(player string, context interface{}) string {
	return "{{b}}pass{{_b}} to spend a chip to pass on this card"
}

func (g *Game) CanPass(player string) bool {
	return g.CurrentlyMoving == player && g.PlayerChips[player] > 0 &&
		!g.IsFinished()
}
