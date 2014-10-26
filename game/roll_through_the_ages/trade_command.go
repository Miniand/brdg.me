package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type TradeCommand struct{}

func (c TradeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("trade", 1, input)
}

func (c TradeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanTrade(pNum)
}

func (c TradeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify how much stone to trade")
	}
	amount, err := strconv.Atoi(a[0])
	if err != nil || amount < 1 {
		return "", errors.New("the amount must be a positive number")
	}

	return "", g.TradeStone(pNum, amount)
}

func (c TradeCommand) Usage(player string, context interface{}) string {
	return "{{b}}trade #{{_b}} to trade stone for 3 workers each, eg. {{b}}trade 3{{_b}}"
}

func (g *Game) CanTrade(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseBuild &&
		g.Boards[player].Developments[DevelopmentEngineering]
}

func (g *Game) TradeStone(player, amount int) error {
	if !g.CanTrade(player) {
		return errors.New("you can't trade at the moment")
	}
	if stone := g.Boards[player].Goods[GoodStone]; amount > stone {
		return fmt.Errorf("you only have %d stone", stone)
	}

	workers := amount * 3
	g.RemainingWorkers += workers
	g.Boards[player].Goods[GoodStone] -= amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s traded {{b}}%d{{_b}} %s for {{b}}%d workers{{_b}}`,
		g.RenderName(player),
		amount,
		RenderGoodName(GoodStone),
		workers,
	)))
	return nil
}
