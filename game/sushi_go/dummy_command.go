package sushi_go

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type DummyCommand struct{}

func (c DummyCommand) Name() string { return "dummy" }

func (c DummyCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify a card to play for the dummy")
	}
	cards := make([]int, len(args))
	for i := range args {
		card, err := strconv.Atoi(args[i])
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
