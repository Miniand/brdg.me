package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
)

type KeepCommand struct{}

func (c KeepCommand) Parse(input string) []string {
	return command.ParseRegexp(`keep`, input)
}

func (c KeepCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.MergerCurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_MERGER
}

func (c KeepCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.KeepShares(playerNum)
}

func (c KeepCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}keep{{_b}} to keep your remaining {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp])
}
