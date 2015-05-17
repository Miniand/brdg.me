package jaipur

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
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
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 2 {
		return "", errors.New("please specify a number and a good type, eg. 3 gold")
	}
	quantity, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("quantity must be a number")
	}
	good, err := helper.MatchStringInStringMap(a[1], GoodStrings)
	if err != nil {
		return "", err
	}
	return "", g.Sell(pNum, quantity, good)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	return "{{b}}sell # [good]{{_b}} to sell goods, eg. {{b}}sell 2 dia{{_b}}"
}

func (g *Game) CanSell(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Sell(player, quantity, good int) error {
	if !g.CanSell(player) {
		return errors.New("you can't sell at the moment")
	}
	minSales, ok := GoodMinSales[good]
	if !ok {
		return errors.New("can't sell that good type")
	}
	if quantity < minSales {
		return fmt.Errorf(
			"you minimum amount you can sell of that good is %d",
			minSales,
		)
	}
	if count := helper.IntCount(good, g.Hands[player]); quantity > count {
		return fmt.Errorf("you only have %d of that good", count)
	}

	numTokens := helper.IntMin(quantity, len(g.Goods[good]))
	g.Tokens[player] = append(g.Tokens[player], g.Goods[good][:numTokens]...)
	g.Goods[good] = g.Goods[good][numTokens:]

	if bonuses, ok := g.Bonuses[quantity]; ok && len(bonuses) > 0 {
		g.Tokens[player] = append(g.Tokens[player], bonuses[0])
		g.Bonuses[quantity] = bonuses[1:]
	}

	g.Hands[player] = helper.IntRemove(good, g.Hands[player], quantity)

	g.NextPlayer()
	return nil
}
