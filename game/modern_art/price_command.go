package modern_art

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type PriceCommand struct{}

func (pc PriceCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("price", 1, input)
}

func (pc PriceCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanSetPrice(player)
}

func (pc PriceCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must specify the price you want to sell for")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	p, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}
	return "", g.SetPrice(playerNum, p)
}

func (pc PriceCommand) Usage(player string, context interface{}) string {
	return "{{b}}price #{{_b}} to set the sale price, eg. {{b}}price 10{{_b}}"
}
