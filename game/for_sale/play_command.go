package for_sale

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (pc PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (pc PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return false
	}
	return g.CanPlay(p)
}

func (pc PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	building, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("building to play must be a number")
	}
	return "", g.Play(p, building)
}

func (pc PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play{{_b}} # to play a building card for sale, eg. {{b}}play 15{{_b}}"
}
