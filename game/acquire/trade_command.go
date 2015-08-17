package acquire

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
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
	if !g.CanTrade(pNum) {
		return "", errors.New("can't trade at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify how many shares to trade")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	return "", g.TradeShares(pNum, g.MergerFromCorp, g.MergerIntoCorp, amount)
}

func (c TradeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}trade #{{_b}} to trade a certain number of your {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares for {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares, two for one.  Eg. {{b}}trade 4{{_b}}`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp],
		CorpColours[g.MergerIntoCorp], CorpNames[g.MergerIntoCorp])
}

func (g *Game) CanTrade(player int) bool {
	return !g.IsFinished() && g.MergerCurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_MERGER &&
		g.PlayerShares[player][g.MergerFromCorp] > 1
}
