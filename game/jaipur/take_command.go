package jaipur

import (
	"errors"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type TakeCommand struct{}

func (c TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("take", 1, -1, input)
}

func (c TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanTake(pNum)
}

func (c TakeCommand) Call(
	player string,
	context interface{},
	args []string,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	takeGoods := []int{}
	forGoods := []int{}
	afterFor := false
	for _, a := range command.ExtractNamedCommandArgs(args) {
		if !afterFor && strings.ToLower(a) == "for" {
			afterFor = true
			continue
		}
		good, err := helper.MatchStringInStringMap(a, GoodStrings)
		if err != nil {
			return "", err
		}
		if afterFor {
			forGoods = append(forGoods, good)
		} else {
			takeGoods = append(takeGoods, good)
		}
	}
	return "", g.Take(pNum, takeGoods, forGoods)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take [goods] (for [goods]){{_b}} to take cards from the market, eg. {{b}}take dia{{_b}} or {{b}}take dia silv for camel spi{{_b}}"
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Take(player int, takeGoods, forGoods []int) error {
	return errors.New("not implemented")
}
