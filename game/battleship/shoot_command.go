package battleship

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type ShootCommand struct{}

func (sc ShootCommand) Name() string { return "shoot" }

func (sc ShootCommand) Call(
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
		return "", errors.New("You must specify the location to shoot")
	}
	y, x, err := ParseLocation(args[0])
	if err != nil {
		return "", err
	}
	return "", g.Shoot(playerNum, y, x)
}

func (sc ShootCommand) Usage(player string, context interface{}) string {
	return "{{b}}shoot ##{{_b}} to shoot at the other player, eg. {{b}}shoot c2{{_b}}"
}
