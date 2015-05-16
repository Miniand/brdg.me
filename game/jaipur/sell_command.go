package jaipur

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type SellCommand struct{}

func (c SellCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("sell", 1, -1, input)
}

func (c SellCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanSell(pNum)
}

func (c SellCommand) Call(
	player string,
	context interface{},
	args []string,
) (string, error) {
	g := context.(*Game)
	_, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	return a[0], errors.New("not implemented")
}

func (c SellCommand) Usage(player string, context interface{}) string {
	return "{{b}}sell # [good]{{_b}} to sell goods, eg. {{b}}sell 2 dia{{_b}}"
}

func (g *Game) CanSell(player int) bool {
	return g.CurrentPlayer == player
}
