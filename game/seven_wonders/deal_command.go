package seven_wonders

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type DealCommand struct{}

func (c DealCommand) Name() string { return "deal" }

func (c DealCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 1 {
		return "", errors.New("you must specify the numbered deal you want to choose")
	}
	dealNum, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("you must specify the numbered deal you want to choose")
	}
	return "", g.Deal(pNum, dealNum-1)
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
	if !g.CanDeal(player) {
		return errors.New("can't choose a deal right now")
	}
	return g.Actions[player].(DealOptioner).ChooseDeal(player, g, deal)
}
