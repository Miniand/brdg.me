package starship_catan

import "github.com/Miniand/brdg.me/command"

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("buy", 1, 2, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	canBuy, _, _ := g.CanBuy(p, ResourceAny)
	return canBuy
}

func (c BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	return "", g.HandleTradeCommand(
		player,
		command.ExtractNamedCommandArgs(args),
		TradeDirBuy,
	)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy #{{_b}} to buy goods, eg. {{b}}buy 3{{_b}}.  If you get to choose which resource to buy you must specify the resource, eg. {{b}}buy 3 food{{_b}}."
}
