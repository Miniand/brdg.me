package starship_catan

import "github.com/Miniand/brdg.me/command"

type FoundColonyCommand struct{}

func (c FoundColonyCommand) Parse(input string) []string {
	return command.ParseNamedCommand("found", input)
}

func (c FoundColonyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFoundColony(p)
}

func (c FoundColonyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.FoundColony(p)
}

func (c FoundColonyCommand) Usage(player string, context interface{}) string {
	return "{{b}}found{{_b}} to found a colony here"
}
