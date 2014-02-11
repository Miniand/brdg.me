package battleship

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type PlaceCommand struct{}

func (pc PlaceCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("place", 3, input)
}

func (pc PlaceCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanPlace(player)
}

func (pc PlaceCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 3 {
		return "", errors.New("You must specify the ship, location, and direction of the ship, eg. place cru b4 down")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	s, err := ParseShip(a[0])
	if err != nil {
		return "", err
	}
	y, x, err := ParseLocation(a[1])
	if err != nil {
		return "", err
	}
	d, err := ParseDirection(a[2])
	if err != nil {
		return "", err
	}
	return "", g.PlaceShip(playerNum, s, y, x, d)
}

func (pc PlaceCommand) Usage(player string, context interface{}) string {
	return "{{b}}place ## ## ##{{_b}} to place a ship on your board, eg. {{b}}place des c2 right{{_b}}"
}
