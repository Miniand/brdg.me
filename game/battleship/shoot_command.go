package battleship

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type ShootCommand struct{}

func (sc ShootCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("shoot", 1, input)
}

func (sc ShootCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanShoot(player)
}

func (sc ShootCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must specify the location to shoot")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	y, x, err := ParseLocation(a[1])
	if err != nil {
		return "", err
	}
	return "", g.Shoot(playerNum, y, x)
}

func (sc ShootCommand) Usage(player string, context interface{}) string {
	return "{{b}}shoot ##{{_b}} to shoot at the other player, eg. {{b}}shoot c2{{_b}}"
}
