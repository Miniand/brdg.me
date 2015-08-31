package acquire

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type KeepCommand struct{}

func (c KeepCommand) Name() string { return "keep" }

func (c KeepCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanKeep(pNum) {
		return "", errors.New("can't keep at the moment")
	}
	return "", g.KeepShares(pNum)
}

func (c KeepCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	return fmt.Sprintf(
		`{{b}}keep{{_b}} to keep your remaining {{b}}{{c "%s"}}%s{{_c}}{{_b}} shares`,
		CorpColours[g.MergerFromCorp], CorpNames[g.MergerFromCorp])
}

func (g *Game) CanKeep(player int) bool {
	return !g.IsFinished() && g.MergerCurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_MERGER
}
