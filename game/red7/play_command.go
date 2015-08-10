package red7

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanPlay(pNum)
}

func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 1 {
		return "", errors.New("you must specify one card")
	}
	card, ok := ParseCard(a[0])
	if !ok {
		return "", errors.New("the card must be a letter followed by a number, eg. r6")
	}

	return "", g.Play(pNum, card)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a card to your palette, eg. {{b}}play b4{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Play(player, card int) error {
	if !g.CanPlay(player) {
		return errors.New("you can't play at the moment")
	}
	return errors.New("not implemented")
}
