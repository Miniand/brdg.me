package modern_art

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
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("You must specify the amount to bid")
	}
	bid, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	return "", g.Bid(playerNum, bid)
}

func (bc BidCommand) Usage(player string, context interface{}) string {
	return "{{b}}bid #{{_b}} to bid an amount in the auction, eg. {{b}}bid 10{{_b}}"
}
