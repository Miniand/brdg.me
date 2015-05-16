package sushizock

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type TakeCommand struct{}

func (tc TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 1, input)
}

func (tc TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanTake(pNum)
}

func (tc TakeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	color, err := helper.MatchStringInStrings(a[0], []string{"blue", "red"})
	if err != nil {
		return "", errors.New(`you must specify "blue" or "red" after "take"`)
	}
	if color == 0 {
		return "", g.TakeBlue(pNum)
	}
	return "", g.TakeRed(pNum)
}

func (tc TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take blue{{_b}} or {{b}}take red{{_b}} to take a blue or red tile"
}
