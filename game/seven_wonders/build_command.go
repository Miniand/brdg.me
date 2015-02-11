package seven_wonders

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type BuildCommand struct{}

func (c BuildCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("build", 1, input)
}

func (c BuildCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanBuild(pNum)
}

func (c BuildCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	_, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to build")
	}
	// TODO finish
	return "", nil
}

func (c BuildCommand) Usage(player string, context interface{}) string {
	return "{{b}}build #{{_b}} to build a card, paying the cost or getting it for free if you have a prerequisite card, eg. {{b}}build 3{{_b}}"
}

func (g *Game) CanBuild(player int) bool {
	return true
}

func (g *Game) CanBuildCard(player int, carder Carder) (
	can bool, coins map[int]map[int]int) {
	c := carder.GetCard()
	for _, freeWith := range c.FreeWith {
		if g.Cards[player].Contains(Cards[freeWith]) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (g *Game) Build(player int, card Carder) error {
	return nil
}
