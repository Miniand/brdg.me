package liars_dice

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strconv"
)

type BidCommand struct{}

func (c BidCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("bid", 2, input)
}

func (c BidCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return !g.IsFinished() && g.WhoseTurn()[0] == player
}

func (c BidCommand) Call(player string, context interface{}, args []string) (
	output string, err error) {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	quantity, err := strconv.Atoi(a[0])
	if err != nil || quantity < 1 {
		return "", errors.New("Quantity must be a positive number, eg. 5")
	}
	if quantity < g.BidQuantity {
		return "", errors.New(fmt.Sprintf(
			"You can't reduce the quantity of the bid, it is currently at %d",
			g.BidQuantity))
	}
	value, err := strconv.Atoi(a[1])
	if err != nil || value < 1 || value > 6 {
		return "", errors.New("Value must be a number between 1 and 6")
	}
	if quantity == g.BidQuantity && value <= g.BidValue {
		return "", errors.New(
			"If you don't increase the bid quantity, you must increase the bid value")
	}
	g.BidQuantity = quantity
	g.BidValue = value
	g.BidPlayer = g.CurrentPlayer
	g.CurrentPlayer = g.NextActivePlayer(g.CurrentPlayer)
	return
}

func (c BidCommand) Usage(player string, context interface{}) string {
	return "{{b}}bid # #{{_b}} to bid the total number of a certain value of dice on the board.  Eg. to bid that there are 8 {{l}}⚄{{_l}} in total, send {{b}}bid 8 5{{_b}}.  {{l}}⚀{{_l}} counts as a wild card for other value dice."
}
