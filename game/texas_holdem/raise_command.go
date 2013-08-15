package texas_holdem

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type RaiseCommand struct{}

func (rc RaiseCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("raise", 1, input)
}

func (rc RaiseCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	currentBet := g.CurrentBet()
	minRaise := g.LargestRaise
	return g.CurrentPlayer == playerNum &&
		g.PlayerMoney[playerNum] > currentBet-g.Bets[playerNum]+minRaise &&
		!g.IsFinished()
}

func (rc RaiseCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must specify the amount to raise by")
	}
	amount, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("Could not understand your raise amount, only use numbers and no punctuation or symbols")
	}
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Raise(playerNum, amount)
}

func (rc RaiseCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	minRaise := g.LargestRaise
	return fmt.Sprintf(
		"{{b}}raise #{{_b}} to raise above the current by the amount, must raise by at least %d, eg. {{b}}raise %d{{_b}}",
		minRaise, minRaise)
}
