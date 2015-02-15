package seven_wonders

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type FreeCommand struct{}

func (c FreeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("free", 1, input)
}

func (c FreeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanFreeBuild(pNum)
}

func (c FreeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to build")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil || cardNum < 1 || cardNum > len(g.Hands[pNum]) {
		return "", errors.New("that is not a valid card number")
	}
	return "", g.FreeBuild(pNum, cardNum-1)
}

func (c FreeCommand) Usage(player string, context interface{}) string {
	return "{{b}}free #{{_b}} to build a card for free, eg. {{b}}free 3{{_b}}"
}

func (g *Game) CanFreeBuild(player int) bool {
	if !g.CanAction(player) {
		return false
	}
	for _, c := range g.Cards[player] {
		if free, ok := c.(FreeBuilder); ok {
			if free.CanFreeBuild() {
				return true
			}
		}
	}
	return false
}

func (g *Game) FreeBuild(player, cardNum int) error {
	if !g.CanFreeBuild(player) {
		return errors.New("cannot free build that card")
	}
	action := &BuildAction{
		Card:   cardNum,
		Free:   true,
		Chosen: true,
	}
	g.Actions[player] = action
	g.CheckHandComplete()
	return nil
}
