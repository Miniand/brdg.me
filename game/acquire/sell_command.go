package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type SellCommand struct{}

func (c SellCommand) Parse(input string) []string {
	return command.ParseRegexp(`sell (\d+)`, input)
}

func (c SellCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.MergerCurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_MERGER &&
		g.PlayerShares[playerNum][g.MergerFromCorp] > 0
}

func (c SellCommand) Call(player string, context interface{},
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
	return "", g.SellShares(playerNum, g.MergerFromCorp, amount)
}

func (c SellCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}sell #{{_b}} to sell a certain number of your {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares.  Eg. {{b}}sell 3{{_b}}`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp])
}
