package sushi_go

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type DummyCommand struct{}

func (c DummyCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("dummy", 1, input)
}

func (c DummyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanDummy(pNum)
}

func (c DummyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	a := command.ExtractNamedCommandArgs(args)
	cards := make([]int, len(a))
	for i := range a {
		card, err := strconv.Atoi(a[i])
		if err != nil {
			return "", errors.New("each card must be a number")
		}
		cards[i] = card - 1 // Input is 1 based
	}

	return "", g.Dummy(pNum, cards)
}

func (c DummyCommand) Usage(player string, context interface{}) string {
	return "{{b}}dummy #{{_b}} to play a card for the dummy, eg. {{b}}dummy 2{{_b}}"
}

func (g *Game) CanDummy(player int) bool {
	return len(g.Players) == 2 && g.Controller == player &&
		g.Playing[Dummy] == nil
}

func (g *Game) Dummy(player int, cards []int) error {
	if !g.CanDummy(player) {
		return errors.New("you can't dummy at the moment")
	}
	l := len(cards)
	if l != 1 {
		return errors.New("you must specify one card to play for the dummy")
	}

	return g.PlayCards(Dummy, player, cards)
}
