package modern_art

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (pc PlayCommand) Name() string { return "play" }

func (pc PlayCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("You must the number of a card to play, such as lmop")
	}
	c, err := ParseCard(args[0])
	if err != nil {
		return "", err
	}
	return "", g.PlayCard(playerNum, c)
}

func (pc PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ####{{_b}} to play a card using the card code, eg. {{b}}play lmop{{_b}}"
}
