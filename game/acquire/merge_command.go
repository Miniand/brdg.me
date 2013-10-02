package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
)

type MergeCommand struct{}

func (c MergeCommand) Parse(input string) []string {
	return command.ParseRegexp(`merge (ARG) into (ARG)`, input)
}

func (c MergeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.CurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_MERGER_CHOOSE
}

func (c MergeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	from, err := CorpFromShortName(args[1])
	if err != nil {
		return "", err
	}
	into, err := CorpFromShortName(args[2])
	if err != nil {
		return "", err
	}
	return "", g.ChooseMerger(from, into)
}

func (c MergeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	availableCommands := []string{}
	for _, merger := range g.PotentialMergers(g.PlayedTile) {
		from := merger[0]
		into := merger[1]
		availableCommands = append(availableCommands,
			fmt.Sprintf(
				`     {{b}}merge %s into %s{{_b}} to merge {{b}}{{c "%s"}}%s{{_c}}{{_b}} into {{b}}{{c "%s"}}%s{{_c}}{{_b}}`,
				CorpShortNames[from], CorpShortNames[into],
				CorpColours[from], CorpNames[from],
				CorpColours[into], CorpNames[into]))
	}
	return `{{b}}merge ## into ##{{_b}} to choose which corporation to merge into another.  Your available options are:`
}
