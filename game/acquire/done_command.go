package acquire

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (c DoneCommand) Name() string { return "done" }

func (c DoneCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanDone(pNum) {
		return "", errors.New("can't call done at the moment")
	}
	g.NextPlayer()
	return "", nil
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return `{{b}}done{{_b}} to end your turn without buying more shares`
}

func (g *Game) CanDone(player int) bool {
	return !g.IsFinished() && g.CurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_BUY_SHARES
}
