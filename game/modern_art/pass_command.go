package modern_art

import "github.com/Miniand/brdg.me/command"

type PassCommand struct{}

func (pc PassCommand) Name() string { return "pass" }

func (pc PassCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	return "", g.Pass(playerNum)
}

func (pc PassCommand) Usage(player string, context interface{}) string {
	return "{{b}}pass{{_b}} to pass"
}
