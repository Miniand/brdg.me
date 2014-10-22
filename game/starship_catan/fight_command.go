package starship_catan

import "github.com/Miniand/brdg.me/command"

type FightCommand struct{}

func (c FightCommand) Parse(input string) []string {
	return command.ParseNamedCommand("fight", input)
}

func (c FightCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFight(p)
}

func (c FightCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Fight(p)
}

func (c FightCommand) Usage(player string, context interface{}) string {
	return "{{b}}fight{{_b}} to fight the pirate"
}
