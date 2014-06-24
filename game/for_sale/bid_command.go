package for_sale

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type BidCommand struct{}

func (bc BidCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("bid", 1, input)
}

func (bc BidCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return false
	}
	return g.CanBid(p)
}

func (bc BidCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	bid, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("bid must be a number")
	}
	return "", g.Bid(p, bid)
}

func (bc BidCommand) Usage(player string, context interface{}) string {
	return "{{b}}bid{{_b}} # to bid a number of your chips to buy a building"
}
