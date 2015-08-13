package jaipur

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type SellCommand struct{}

func (c SellCommand) Name() string { return "sell" }

func (c SellCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify what to sell, eg. 5 go or dia dia")
	}
	goods, err := ParseGoodStr(strings.Join(args, " "))
	if err != nil {
		return "", err
	}
	if len(goods) == 0 {
		return "", errors.New("please specify what to sell, eg. 5 go or dia dia")
	}
	quantity := 0
	good := goods[0]
	if _, ok := helper.IntFind(good, TradeGoods); !ok {
		return "", errors.New("please specify a trade good")
	}
	for _, g := range goods {
		if g != good {
			return "", errors.New("please only specify one type of good to sell")
		}
		quantity++
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
	points := helper.IntSum(g.Goods[good][:numTokens])
	g.Tokens[player] = append(g.Tokens[player], g.Goods[good][:numTokens]...)
	g.Goods[good] = g.Goods[good][numTokens:]
	g.GoodTokens[player] += numTokens

	suffix := ""
	if bonuses, ok := g.Bonuses[quantity]; ok && len(bonuses) > 0 {
		g.Tokens[player] = append(g.Tokens[player], bonuses[0])
		g.BonusTokens[player]++
		g.Bonuses[quantity] = bonuses[1:]
		suffix = " and took a bonus token"
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s sold {{b}}%d %s{{_b}} for {{b}}%d %s{{_b}}%s",
		g.RenderName(player),
		quantity,
		render.Colour(helper.Plural(quantity, GoodStrings[good]), GoodColours[good]),
		points,
		helper.Plural(points, "point"),
		suffix,
	)))
	if suffix != "" {
		g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
			"The bonus token was {{b}}%d points{{_b}}",
			g.Tokens[player][len(g.Tokens[player])-1],
		), []string{g.Players[player]}))
	}

	g.Hands[player] = helper.IntRemove(good, g.Hands[player], quantity)

	// End of round if three good types are out of tokens.
	emptiedGoods := 0
	for _, good := range TradeGoods {
		if len(g.Goods[good]) == 0 {
			emptiedGoods++
		}
		if emptiedGoods >= 3 {
			g.EndRound()
			return nil
		}
	}
	g.NextPlayer()
	return nil
}
