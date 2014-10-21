package starship_catan

import "github.com/Miniand/brdg.me/command"

type SellCommand struct{}

func (c SellCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("sell", 1, 2, input)
}

func (c SellCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	canSell, _, _ := g.CanSell(p, ResourceAny)
	return canSell
}

func (c SellCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	return "", g.HandleTradeCommand(
		player,
		command.ExtractNamedCommandArgs(args),
		TradeDirSell,
	)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	return "{{b}}sell #{{_b}} to sell goods, eg. {{b}}sell 3{{_b}}.  If you get to choose which resource to sell you must specify the resource, eg. {{b}}sell 3 food{{_b}}."
}
