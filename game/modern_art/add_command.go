package modern_art

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type AddCommand struct{}

func (ac AddCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("add", 1, input)
}

func (ac AddCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanAdd(player)
}

func (ac AddCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must the number of a card to play, such as 2")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	c, err := ParseCard(a[0])
	if err != nil {
		return "", err
	}
	return "", g.AddCard(playerNum, c)
}

func (ac AddCommand) Usage(player string, context interface{}) string {
	return "{{b}}play #{{_b}} to play a card, eg. {{b}}play 2{{_b}}"
}
