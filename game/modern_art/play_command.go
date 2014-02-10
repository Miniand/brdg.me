package modern_art

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (pc PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (pc PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanPlay(player)
}

func (pc PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must the number of a card to play, such as lmop")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	c, err := ParseCard(a[0])
	if err != nil {
		return "", err
	}
	return "", g.PlayCard(playerNum, c)
}

func (pc PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ####{{_b}} to play a card using the card code, eg. {{b}}play lmop{{_b}}"
}
