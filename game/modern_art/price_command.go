package modern_art

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type PriceCommand struct{}

func (pc PriceCommand) Name() string { return "price" }

func (pc PriceCommand) Call(
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
		return "", errors.New("You must specify the price you want to sell for")
	}
	p, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	return "", g.SetPrice(playerNum, p)
}

func (pc PriceCommand) Usage(player string, context interface{}) string {
	return "{{b}}price #{{_b}} to set the sale price, eg. {{b}}price 10{{_b}}"
}
