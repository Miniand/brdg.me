package starship_catan

import "github.com/Miniand/brdg.me/command"

type FoundTradeCommand struct{}

func (c FoundTradeCommand) Parse(input string) []string {
	return command.ParseNamedCommand("found", input)
}

func (c FoundTradeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFoundTradingPost(p)
}

func (c FoundTradeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.FoundTradingPost(p)
}

func (c FoundTradeCommand) Usage(player string, context interface{}) string {
	return "{{b}}found{{_b}} to found a trading post here"
}
