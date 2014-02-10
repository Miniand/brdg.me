package modern_art

import (
	"github.com/Miniand/brdg.me/command"
)

type BuyCommand struct{}

func (bc BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommand("buy", input)
}

func (bc BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.CanBuy(player)
}

func (bc BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	return "", g.Buy(playerNum)
}

func (bc BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy{{_b}} to buy the current card for the set price"
}
