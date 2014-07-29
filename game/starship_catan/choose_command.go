package starship_catan

import "github.com/Miniand/brdg.me/command"

type ChooseCommand struct{}

func (c ChooseCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("choose", 1, input)
}

func (c ChooseCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanChoose(p)
}

func (c ChooseCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	m, err := ParseModule(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Choose(p, m)
}

func (c ChooseCommand) Usage(player string, context interface{}) string {
	return "{{b}}choose ##{{_b}} to choose which module to start with.  The logistics module is the most useful module for starting players.  Eg. {{b}}choose lo{{_b}}"
}
