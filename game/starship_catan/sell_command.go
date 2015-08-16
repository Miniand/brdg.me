package starship_catan

import "github.com/Miniand/brdg.me/command"

type SellCommand struct{}

func (c SellCommand) Name() string { return "sell" }

func (c SellCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	return "", g.HandleTradeCommand(
		player,
		input,
		TradeDirSell,
	)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	return "{{b}}sell #{{_b}} to sell goods, eg. {{b}}sell 3{{_b}}.  If you get to choose which resource to sell you must specify the resource, eg. {{b}}sell 3 food{{_b}}."
}

func (g *Game) CanSell(player int) bool {
	ok, _, _ := g.CanSellResource(player, ResourceAny)
	return ok
}

func (g *Game) CanSellResource(player, resource int) (ok bool, price int, reason string) {
	return g.CanTrade(player, resource, TradeDirSell)
}

type TradePhaseSellCommand struct {
	SellCommand
}

func (g *Game) CanTradePhaseSell(player int) bool {
	return g.Phase == PhaseTradeAndBuild && g.CanSell(player)
}
