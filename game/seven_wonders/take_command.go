package seven_wonders

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type TakeCommand struct{}

func (c TakeCommand) Name() string { return "take" }

func (c TakeCommand) Call(
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
		return "", errors.New("you must specify which numbered card to take")
	}
	cardNum, err := strconv.Atoi(args[0])
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

func (g *Game) CanTakeCard(player int, carder Carder) bool {
	c := carder.GetCard()
	// See if you already have it, which means you can't get it again.
	for _, pc := range g.Cards[player] {
		if pc.(Carder).GetCard().Name == c.Name {
			return false
		}
	}
	return true
}

func (g *Game) Take(player, cardNum int) error {
	if !g.CanTake(player) {
		return errors.New("cannot take at the moment")
	}
	if cardNum < 0 || cardNum >= len(g.Discard) {
		return errors.New("invalid card number")
	}

	c := g.Discard[cardNum]
	if !g.CanTakeCard(player, c.(Carder)) {
		return errors.New("you can't have more than one of the same card")
	}
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
