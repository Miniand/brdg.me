package texas_holdem

import (
	"github.com/Miniand/brdg.me/command"
)

type AllinCommand struct{}

func (ac AllinCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("allin", 0, input)
}

func (ac AllinCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CurrentPlayer == playerNum && g.PlayerMoney[playerNum] > 0 &&
		!g.IsFinished()
}

func (ac AllinCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return err
	}
	return g.AllIn(playerNum)
}

func (ac AllinCommand) Usage(player string, context interface{}) string {
	return "{{b}}allin{{_b}} to bet all your money and go all in"
}
