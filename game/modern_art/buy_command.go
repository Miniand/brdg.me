package modern_art

import "github.com/Miniand/brdg.me/command"

type BuyCommand struct{}

func (bc BuyCommand) Name() string { return "buy" }

func (bc BuyCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
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
