package acquire

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
)

type MergeCommand struct{}

func (c MergeCommand) Name() string { return "merge" }

func (c MergeCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanMerge(pNum) {
		return "", errors.New("can't merge at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 3 || strings.ToLower(args[1]) != "into" {
		return "", errors.New("you must specify which hotel to merge, eg. 'merge fe into im'")
	}
	from, err := FindCorp(args[0])
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

func (g *Game) CanMerge(player int) bool {
	return !g.IsFinished() && g.CurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_MERGER_CHOOSE
}
