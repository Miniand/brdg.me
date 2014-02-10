package modern_art

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type BidCommand struct{}

func (bc BidCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("bid", 1, input)
}

func (bc BidCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanBid(player)
}

func (bc BidCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must specify the amount to bid")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	bid, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Bid(playerNum, bid)
}

func (bc BidCommand) Usage(player string, context interface{}) string {
	return "{{b}}bid #{{_b}} to bid an amount in the auction, eg. {{b}}bid 10{{_b}}"
}
