package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type TradeCommand struct{}

func (c TradeCommand) Name() string { return "trade" }

func (c TradeCommand) Call(
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
	if err != nil || len(args) == 0 {
		return "", errors.New("you must specify how much stone to trade")
	}
	amount, err := strconv.Atoi(args[0])
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
		g.Boards[player].Developments[DevelopmentEngineering] &&
		g.Boards[player].Goods[GoodStone] > 0
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
