package sushizock

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (rc RollCommand) Name() string { return "roll" }

func (rc RollCommand) Call(
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
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify something to roll")
	}
	dice := make([]int, len(args))
	for i, s := range args {
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
