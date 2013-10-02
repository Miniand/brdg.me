package acquire

import (
	"github.com/Miniand/brdg.me/command"
)

type FoundCommand struct{}

func (c FoundCommand) Parse(input string) []string {
	return command.ParseRegexp(`found (ARG)`, input)
}

func (c FoundCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.CurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_FOUND_CORP
}

func (c FoundCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	corp, err := CorpFromShortName(args[1])
	if err != nil {
		return "", err
	}
	return "", g.FoundCorp(playerNum, corp)
}

func (c FoundCommand) Usage(player string, context interface{}) string {
	return `{{b}}found ##{{_b}} to found a corporation on your tile.  Eg. {{b}}found fe{{_b}}`
}
