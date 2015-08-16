package modern_art

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type AddCommand struct{}

func (ac AddCommand) Name() string { return "add" }

func (ac AddCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
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
	return "", g.AddCard(playerNum, c)
}

func (ac AddCommand) Usage(player string, context interface{}) string {
	return "{{b}}add ####{{_b}} to add a card to the auction, eg. {{b}}play lmop{{_b}}"
}
