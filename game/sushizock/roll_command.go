package sushizock

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (rc RollCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("roll", 1, -1, input)
}

func (rc RollCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanRoll(pNum)
}

func (rc RollCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	dice := make([]int, len(a))
	for i, s := range a {
		d, err := strconv.Atoi(s)
		if err != nil {
			return "", fmt.Errorf("%s is not a number", s)
		}
		dice[i] = d
	}
	return "", g.RollDice(pNum, dice)
}

func (rc RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll # # #{{_b}} to reroll dice, eg. {{b}}roll 1 3 4{{_b}}"
}
