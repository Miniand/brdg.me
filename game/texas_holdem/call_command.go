package texas_holdem

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
)

type CallCommand struct{}

func (cc CallCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("call", 0, input)
}

func (cc CallCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == playerNum && g.Bets[playerNum] < currentBet &&
		g.PlayerMoney[playerNum] > currentBet-g.Bets[playerNum] &&
		!g.IsFinished()
}

func (cc CallCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Call(playerNum)
}

func (cc CallCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		"{{b}}call{{_b}} to increase your bet by %d to match the current bet",
		g.CurrentBet()-g.Bets[playerNum])
}