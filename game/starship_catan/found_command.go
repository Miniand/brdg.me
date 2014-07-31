package starship_catan

import "github.com/Miniand/brdg.me/command"

type FoundCommand struct{}

func (c FoundCommand) Parse(input string) []string {
	return command.ParseNamedCommand("found", input)
}

func (c FoundCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFound(p)
}

func (c FoundCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Found(p)
}

func (c FoundCommand) Usage(player string, context interface{}) string {
	return "{{b}}found{{_b}} to found a colony here"
}
