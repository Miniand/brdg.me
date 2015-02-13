package seven_wonders

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DealCommand struct{}

func (c DealCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("deal", 1, input)
}

func (c DealCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanDeal(pNum)
}

func (c DealCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	_, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify the numbered deal")
	}
	// TODO finish
	return "", nil
}

func (c DealCommand) Usage(player string, context interface{}) string {
	return "{{b}}deal #{{_b}} to choose which trade deal you want, eg. {{b}}deal 2{{_b}}"
}

func (g *Game) CanDeal(player int) bool {
	if g.Actions[player] == nil {
		return false
	}
	_, ok := g.Actions[player].(DealOptioner)
	return ok && !g.Actions[player].IsComplete()
}

func (g *Game) Deal(player, deal int) error {
	return nil
}
