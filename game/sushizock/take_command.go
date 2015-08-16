package sushizock

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type TakeCommand struct{}

func (tc TakeCommand) Name() string { return "take" }

func (tc TakeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a color to take")
	}
	color, err := helper.MatchStringInStrings(args[0], []string{"blue", "red"})
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
