package acquire

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type BuyCommand struct{}

func (c BuyCommand) Name() string { return "buy" }

func (c BuyCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanBuy(pNum) {
		return "", errors.New("can't buy at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 2 {
		return "", errors.New("you must specify a number and a hotel")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	corp, err := FindCorp(args[1])
	if err != nil {
		return "", err
	}
	return "", g.BuyShares(pNum, corp, amount)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return `{{b}}buy # ##{{_b}} to buy a certain number of shares in a corp.  Eg. {{b}}buy 3 worldwide{{_b}} or {{b}}buy 3 wo{{_b}}`
}

func (g *Game) CanBuy(player int) bool {
	return !g.IsFinished() && g.CurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_BUY_SHARES
}
