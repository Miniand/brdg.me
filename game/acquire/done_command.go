package acquire

import (
	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (c DoneCommand) Parse(input string) []string {
	return command.ParseRegexp(`done`, input)
}

func (c DoneCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.CurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_BUY_SHARES
}

func (c DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	g.NextPlayer()
	return "", nil
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return `{{b}}done{{_b}} to end your turn without buying more shares`
}
