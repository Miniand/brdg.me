package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type TradeCommand struct{}

func (c TradeCommand) Parse(input string) []string {
	return command.ParseRegexp(`trade (\d+)`, input)
}

func (c TradeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.MergerCurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_MERGER &&
		g.PlayerShares[playerNum][g.MergerFromCorp] > 1
}

func (c TradeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.TradeShares(playerNum, g.MergerFromCorp, g.MergerIntoCorp,
		amount)
}

func (c TradeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}trade #{{_b}} to trade a certain number of your {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares for {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares, two for one.  Eg. {{b}}trade 4{{_b}}`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp],
		CorpColours[g.MergerIntoCorp], CorpNames[g.MergerIntoCorp])
}
