package starship_catan

import "github.com/Miniand/brdg.me/command"

type BuyCommand struct{}

func (c BuyCommand) Name() string { return "buy" }

func (c BuyCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	return "", g.HandleTradeCommand(
		player,
		input,
		TradeDirBuy,
	)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy #{{_b}} to buy goods, eg. {{b}}buy 3{{_b}}.  If you get to choose which resource to buy you must specify the resource, eg. {{b}}buy 3 food{{_b}}."
}

func (g *Game) CanBuy(player int) bool {
	ok, _, _ := g.CanBuyResource(player, ResourceAny)
	return ok
}

func (g *Game) CanBuyResource(player, resource int) (ok bool, price int, reason string) {
	return g.CanTrade(player, resource, TradeDirBuy)
}

type TradePhaseBuyCommand struct {
	BuyCommand
}

func (g *Game) CanTradePhaseBuy(player int) bool {
	return g.Phase == PhaseTradeAndBuild && g.CanBuy(player)
}
