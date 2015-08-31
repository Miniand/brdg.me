package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type SwapCommand struct{}

func (c SwapCommand) Name() string { return "swap" }

func (c SwapCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 3 {
		return "", errors.New("you must specify an amount, a good to remove and a good to receive")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil || amount < 1 {
		return "", errors.New("the amount must be a positive number")
	}

	fromGood, err := helper.MatchStringInStringMap(args[1], GoodStrings)
	if err != nil {
		return "", err
	}

	toGood, err := helper.MatchStringInStringMap(args[2], GoodStrings)
	if err != nil {
		return "", err
	}

	return "", g.Swap(pNum, fromGood, toGood, amount)
}

func (c SwapCommand) Usage(player string, context interface{}) string {
	return "{{b}}swap # (from) (to){{_b}} to swap goods from one type to another, eg. {{b}}swap 2 wood spear{{_b}}"
}

func (g *Game) CanSwap(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseTrade &&
		g.Boards[player].Developments[DevelopmentShipping] &&
		g.Boards[player].GoodsNum() > 0
}

func (g *Game) Swap(player, fromGood, toGood, amount int) error {
	if !g.CanSwap(player) {
		return errors.New("you can't swap at the moment")
	}
	if amount < 1 {
		return errors.New("amount must be positive")
	}
	if fromGood == toGood {
		return errors.New("you must specify two different goods")
	}
	if amount > g.RemainingShips {
		return fmt.Errorf("you only have %d ships remaining", g.RemainingShips)
	}
	if goodNum := g.Boards[player].Goods[fromGood]; goodNum < amount {
		return fmt.Errorf(
			"you only have %d %s left",
			goodNum,
			GoodStrings[fromGood],
		)
	}
	if max := GoodMaximum(toGood); g.Boards[player].Goods[toGood]+amount > max {
		return fmt.Errorf(
			"the you only have room for %d %s",
			max,
			GoodStrings[toGood],
		)
	}

	g.Boards[player].Goods[fromGood] -= amount
	g.Boards[player].Goods[toGood] += amount
	g.RemainingShips -= amount

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s swapped {{b}}%d{{_b}} %s for %s`,
		g.RenderName(player),
		amount,
		RenderGoodName(fromGood),
		RenderGoodName(toGood),
	)))

	if g.RemainingShips == 0 {
		g.NextPhase()
	}
	return nil
}
