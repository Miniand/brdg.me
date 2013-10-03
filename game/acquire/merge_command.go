package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"strings"
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
	from, err := FindCorp(args[1])
	if err != nil {
		return "", err
	}
	into, err := FindCorp(args[2])
	if err != nil {
		return "", err
	}
	return "", g.ChooseMerger(g.PlayedTile, from, into)
}

func (c MergeCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	availableCommands := []string{}
	for _, merger := range g.PotentialMergers(g.PlayedTile) {
		from := merger[0]
		into := merger[1]
		availableCommands = append(availableCommands,
			fmt.Sprintf(
				`     {{b}}merge %s into %s{{_b}} to merge {{b}}%s{{_b}} into {{b}}%s{{_b}}`,
				CorpShortNames[from], CorpShortNames[into],
				RenderCorp(from),
				RenderCorp(into)))
	}
	return fmt.Sprintf(
		"{{b}}merge ## into ##{{_b}} to choose which corporation to merge into another.  Your available options are:\n%s",
		strings.Join(availableCommands, "\n"))
}
