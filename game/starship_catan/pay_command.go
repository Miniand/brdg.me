package starship_catan

import "github.com/Miniand/brdg.me/command"

type PayCommand struct{}

func (c PayCommand) Parse(input string) []string {
	return command.ParseNamedCommand("pay", input)
}

func (c PayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanPayRansom(p)
}

func (c PayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.PayRansom(p)
}

func (c PayCommand) Usage(player string, context interface{}) string {
	return "{{b}}pay{{_b}} to pay the ransom"
}
