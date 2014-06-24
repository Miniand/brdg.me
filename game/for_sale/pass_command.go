package for_sale

import "github.com/Miniand/brdg.me/command"

type PassCommand struct{}

func (pc PassCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("pass", 0, input)
}

func (pc PassCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return false
	}
	return g.CanBid(p)
}

func (pc PassCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Pass(p)
}

func (pc PassCommand) Usage(player string, context interface{}) string {
	return "{{b}}pass{{_b}} to pass and take the lowest building and get half of your bid back (rounded up)"
}
