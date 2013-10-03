package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type EndCommand struct{}

func (c EndCommand) Parse(input string) []string {
	return command.ParseRegexp(`end`, input)
}

func (c EndCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanEnd(playerNum)
}

func (c EndCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	g.FinalTurn = true
	g.Log = g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} triggered the end of the game at the end of their turn`,
		g.RenderPlayer(g.CurrentPlayer))))
	return "", nil
}

func (c EndCommand) Usage(player string, context interface{}) string {
	return `{{b}}end{{_b}} to end the game after your current turn`
}
