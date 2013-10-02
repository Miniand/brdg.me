package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseRegexp(`buy (\d+) (ARG)`, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.CurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_BUY_SHARES
}

func (c BuyCommand) Call(player string, context interface{},
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
	corp, err := CorpFromShortName(args[2])
	if err != nil {
		return "", err
	}
	return "", g.BuyShares(playerNum, corp, amount)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}sell #{{_b}} to sell a certain number of your {{b}}{c "%s"}}%s{_c}}{{_b}} shares.  Eg. {{b}}sell 3{{_b}}`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp])
}
