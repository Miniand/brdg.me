package modern_art

import (
	"github.com/Miniand/brdg.me/command"
)

type PassCommand struct{}

func (pc PassCommand) Parse(input string) []string {
	return command.ParseNamedCommand("pass", input)
}

func (pc PassCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanPass(player)
}

func (pc PassCommand) Call(player string, context interface{},
	args []string) (string, error) {
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
