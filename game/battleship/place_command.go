package battleship

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type PlaceCommand struct{}

func (pc PlaceCommand) Name() string { return "place" }

func (pc PlaceCommand) Call(
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
	if err != nil || len(args) != 3 {
		return "", errors.New("You must specify the ship, location, and direction of the ship, eg. place cru b4 down")
	}
	s, err := ParseShip(args[0])
	if err != nil {
		return "", err
	}
	y, x, err := ParseLocation(args[1])
	if err != nil {
		return "", err
	}
	d, err := ParseDirection(args[2])
	if err != nil {
		return "", err
	}
	return "", g.PlaceShip(playerNum, s, y, x, d)
}

func (pc PlaceCommand) Usage(player string, context interface{}) string {
	return "{{b}}place ## ## ##{{_b}} to place a ship on your board, eg. {{b}}place des c2 right{{_b}}"
}
