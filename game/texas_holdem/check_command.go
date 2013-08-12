package texas_holdem

import (
	"github.com/beefsack/brdg.me/command"
)

type CheckCommand struct{}

func (cc CheckCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("check", 0, input)
}

func (cc CheckCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == playerNum && g.Bets[playerNum] == currentBet &&
		!g.IsFinished()
}

func (cc CheckCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return err
	}
	return g.Check(playerNum)
}

func (cc CheckCommand) Usage(player string, context interface{}) string {
	return "{{b}}check{{_b}} to continue without betting more money"
}
