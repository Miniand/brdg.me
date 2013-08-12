package lost_cities

import (
	"github.com/beefsack/brdg.me/command"
)

type ReadyCommand struct{}

func (d ReadyCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("ready", 0, input)
}

func (d ReadyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return false
	}
	return !g.ReadyPlayers[playerNum] && g.IsEndOfRound()
}

func (d ReadyCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	return g.PlayerReady(playerNum)
}

func (d ReadyCommand) Usage(player string, context interface{}) string {
	return "{{b}}ready{{_b}} to start the next round"
}
