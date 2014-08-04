package sushizock

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type StealCommand struct{}

func (sc StealCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("steal", 2, 3, input)
}

func (sc StealCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanSteal(pNum)
}

func (sc StealCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	target, err := helper.MatchStringInStrings(a[0], g.Players)
	if err != nil {
		return "", err
	}
	color, err := helper.MatchStringInStrings(a[1], []string{"blue", "red"})
	if err != nil {
		return "", errors.New(`you must specify "blue" or "red" after "take"`)
	}
	if len(a) == 3 {
		n, err := strconv.Atoi(a[2])
		if err != nil {
			return "", fmt.Errorf("%s is not a number", a[2])
		}
		if color == 0 {
			return "", g.StealBlueN(pNum, target, n)
		}
		return "", g.StealRedN(pNum, target, n)
	} else {
		if color == 0 {
			return "", g.StealBlue(pNum, target)
		}
		return "", g.StealRed(pNum, target)
	}
}

func (sc StealCommand) Usage(player string, context interface{}) string {
	return "{{b}}steal player blue/red (#){{_b}} to steal a blue or red card from another player, eg {{b}}steal mick red{{_b}}.  If you have four or more chopsticks you can choose which tile in the stack to steal with 1 being the top, eg. {{b}}steal steve blue 2{{_b}}"
}
