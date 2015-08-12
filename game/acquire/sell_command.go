package acquire

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type SellCommand struct{}

func (c SellCommand) Name() string { return "sell" }

func (c SellCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanSell(pNum) {
		return "", errors.New("can't sell at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil {
		return "", errors.New("please specify how many shares to sell")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	return "", g.SellSharesAction(pNum, g.MergerFromCorp, amount)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}sell #{{_b}} to sell a certain number of your {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares.  Eg. {{b}}sell 3{{_b}}`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp])
}

func (g *Game) CanSell(player int) bool {
	return !g.IsFinished() && g.MergerCurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_MERGER &&
		g.PlayerShares[player][g.MergerFromCorp] > 0
}
