package alhambra

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type TakeCommand struct{}

func (c TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("take", 1, -1, input)
}

func (c TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanTake(pNum)
}

func (c TakeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", ErrCouldNotFindPlayer
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which cards to take")
	}
	cards := card.Deck{}
	for _, rawCard := range a {
		c, err := ParseCard(rawCard)
		if err != nil {
			return "", err
		}
		cards = cards.Push(c)
	}
	return "", g.Take(pNum, cards)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take ## (##){{_b}} to take multiple cards up to the value of 5, or a single card over the value of 5, eg. {{b}}take r2 b3{{_b}}"
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseAction
}

func (g *Game) Take(player int, cards card.Deck) error {
	if !g.CanTake(player) {
		return errors.New("unable to take right now")
	}
	if cards.Len() == 0 {
		return errors.New("must specify at least one card to take with")
	}
	total := 0
	for _, c := range cards {
		if g.Cards.Contains(c) == 0 {
			return fmt.Errorf("%s isn't in the available cards", c)
		}
		total += c.(Card).Value
	}
	if cards.Len() > 1 && total > 5 {
		return errors.New("you can't take multiple cards if their combined value is greater than 5")
	}

	// Take the cards
	g.Boards[player].Cards = g.Boards[player].Cards.PushMany(cards)
	for _, c := range cards {
		g.Cards, _ = g.Cards.Remove(c, 1)
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %s",
		g.PlayerName(player),
		RenderCards(cards),
	)))

	g.NextPhase()
	return nil
}
