package texas_holdem

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type RaiseCommand struct{}

func (rc RaiseCommand) Name() string { return "raise" }

func (rc RaiseCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 1 {
		return "", errors.New("You must specify the amount to raise by")
	}
	amount, err := strconv.Atoi(args[0])
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

func (g *Game) CanRaise(player int) bool {
	currentBet := g.CurrentBet()
	minRaise := g.LargestRaise
	return g.CurrentPlayer == player &&
		g.PlayerMoney[player] > currentBet-g.Bets[player]+minRaise &&
		!g.IsFinished()
}
