package for_sale

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (pc PlayCommand) Name() string { return "play" }

func (pc PlayCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a building")
	}
	building, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("building to play must be a number")
	}
	return "", g.Play(p, building)
}

func (pc PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play{{_b}} # to play a building card for sale, eg. {{b}}play 15{{_b}}"
}
