package for_sale

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type BidCommand struct{}

func (bc BidCommand) Name() string { return "bid" }

func (bc BidCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a bid")
	}
	bid, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("bid must be a number")
	}
	return "", g.Bid(p, bid)
}

func (bc BidCommand) Usage(player string, context interface{}) string {
	return "{{b}}bid{{_b}} # to bid a number of your chips to buy a building"
}
