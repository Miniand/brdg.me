package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type DiscardCommand struct{}

func (c DiscardCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("discard", 2, input)
}

func (c DiscardCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanDiscard(pNum)
}

func (c DiscardCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 2 {
		return "", errors.New(
			"you must pass an amount to discard and the name of a thing to discard")
	}

	amount, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("you must specify an amount")
	}

	good, err := helper.MatchStringInStringMap(a[1], GoodStrings)
	if err != nil {
		return "", err
	}

	return "", g.Discard(pNum, amount, good)
}

func (c DiscardCommand) Usage(player string, context interface{}) string {
	return "{{b}}discard # (good){{_b}} to discard goods down to the required 6, eg. {{b}}discard 2 wood{{_b}}"
}

func (g *Game) CanDiscard(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseDiscard
}

func (g *Game) Discard(player, amount, good int) error {
	if !g.CanDiscard(player) {
		return errors.New("you can't discard at the moment")
	}
	if amount < 1 {
		return errors.New("amount must be a positive number")
	}
	if !ContainsInt(good, Goods) {
		return errors.New("invalid good")
	}
	if num := g.Boards[player].Goods[good]; amount > num {
		return fmt.Errorf("you only have %d %s", num, GoodStrings[good])
	}
	if num := g.Boards[player].GoodsNum(); num-amount < 6 {
		return fmt.Errorf("you only need to discard %d", num-6)
	}
	g.Boards[player].Goods[good] -= amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %d %s",
		g.RenderName(player),
		amount,
		RenderGoodName(good),
	)))
	if g.Boards[player].GoodsNum() <= 6 {
		g.NextTurn()
	}
	return nil
}
