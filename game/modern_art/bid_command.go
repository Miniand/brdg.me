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
		return "", errors.New("You must the number of a card to play, such as 2")
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
	return "{{b}}play #{{_b}} to play a card, eg. {{b}}play 2{{_b}}"
}
