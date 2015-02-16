package seven_wonders

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type WonderCommand struct{}

func (c WonderCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("wonder", 1, input)
}

func (c WonderCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanWonder(pNum)
}

func (c WonderCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to use to build the wonder stage")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil || cardNum < 1 || cardNum > len(g.Hands[pNum]) {
		return "", errors.New("that is not a valid card number")
	}
	return "", g.Wonder(pNum, cardNum-1)
}

func (c WonderCommand) Usage(player string, context interface{}) string {
	return "{{b}}wonder #{{_b}} to build a wonder stage using the specified card, eg. {{b}}wonder 3{{_b}}"
}

func (g *Game) CanWonder(player int) bool {
	remaining := g.RemainingWonderStages(player)
	if !g.CanAction(player) || remaining.Len() == 0 {
		return false
	}
	can, _ := g.CanAfford(player, remaining[0].(Carder).GetCard().Cost)
	return can
}

func (g *Game) Wonder(player, cardNum int) error {
	if !g.CanWonder(player) {
		return errors.New("cannot build a wonder stage right now")
	}
	_, coins := g.CanAfford(player,
		g.RemainingWonderStages(player)[0].(Carder).GetCard().Cost)
	action := &BuildAction{
		Card:   cardNum,
		Wonder: true,
	}
	if len(coins) <= 1 {
		action.Chosen = true
	}
	g.Actions[player] = action
	g.CheckHandComplete()
	return nil
}
