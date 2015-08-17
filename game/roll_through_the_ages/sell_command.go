package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type SellCommand struct{}

func (c SellCommand) Name() string { return "sell" }

func (c SellCommand) Call(
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
		return "", errors.New("you must specify how much food to sell")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil || amount < 1 {
		return "", errors.New("the amount must be a positive number")
	}

	return "", g.SellFood(pNum, amount)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	return "{{b}}sell #{{_b}} to sell food for 6 coins each, eg. {{b}}sell 3{{_b}}"
}

func (g *Game) CanSell(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseBuy &&
		g.Boards[player].Developments[DevelopmentGranaries]
}

func (g *Game) SellFood(player, amount int) error {
	if !g.CanSell(player) {
		return errors.New("you can't sell at the moment")
	}
	if amount > g.Boards[player].Food {
		return fmt.Errorf("you only have %d food", g.Boards[player].Food)
	}

	coins := amount * 6
	g.RemainingCoins += coins
	g.Boards[player].Food -= amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s traded {{b}}%d{{_b}} %s for {{b}}%d coins{{_b}}`,
		g.RenderName(player),
		amount,
		FoodName,
		coins,
	)))
	return nil
}
