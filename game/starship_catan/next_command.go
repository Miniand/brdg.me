package starship_catan

import "github.com/Miniand/brdg.me/command"

type NextCommand struct{}

func (c NextCommand) Parse(input string) []string {
	return command.ParseNamedCommand("next", input)
}

func (c NextCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanNext(p)
}

func (c NextCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Next(p)
}

func (c NextCommand) Usage(player string, context interface{}) string {
	return "{{b}}next{{_b}} to advance to the next card"
}
