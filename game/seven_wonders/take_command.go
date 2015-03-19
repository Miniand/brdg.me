package seven_wonders

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type TakeCommand struct{}

func (c TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 1, input)
}

func (c TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanTake(pNum)
}

func (c TakeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which numbered card to take")
	}
	cardNum, err := strconv.Atoi(a[0])
	if err != nil || cardNum < 1 || cardNum > len(g.Discard) {
		return "", errors.New("that is not a valid card number")
	}
	return "", g.Take(pNum, cardNum-1)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take #{{_b}} to take a card from the discard pile, eg. {{b}}take 3{{_b}}"
}

func (g *Game) CanTake(player int) bool {
	return len(g.ToResolve) > 0 && InStrs(
		g.Players[player],
		g.ToResolve[0].WhoseTurn(g),
	)
}

func (g *Game) Take(player, cardNum int) error {
	if !g.CanTake(player) {
		return errors.New("cannot take at the moment")
	}
	if cardNum < 0 || cardNum >= len(g.Discard) {
		return errors.New("invalid card number")
	}

	c := g.Discard[cardNum]
	g.Cards[player] = g.Cards[player].Push(c)
	g.Discard = append(g.Discard[:cardNum], g.Discard[cardNum+1:]...)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %s",
		g.PlayerName(player),
		RenderCard(c.(Carder)),
	)))
	if handler, ok := c.(PostActionExecuteHandler); ok {
		handler.HandlePostActionExecute(player, g)
	}

	g.Resolved()

	return nil
}
