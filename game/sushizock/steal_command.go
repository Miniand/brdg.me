package sushizock

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type StealCommand struct{}

func (sc StealCommand) Name() string { return "steal" }

func (sc StealCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 2 {
		return "", errors.New("you must specify at least a player and a color")
	}
	target, err := helper.MatchStringInStrings(args[0], g.Players)
	if err != nil {
		return "", err
	}
	color, err := helper.MatchStringInStrings(args[1], []string{"blue", "red"})
	if err != nil {
		return "", errors.New(`you must specify "blue" or "red" after the player name`)
	}
	if len(args) == 3 {
		n, err := strconv.Atoi(args[2])
		if err != nil {
			return "", fmt.Errorf("%s is not a number", args[2])
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
