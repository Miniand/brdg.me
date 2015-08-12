package acquire

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type FoundCommand struct{}

func (c FoundCommand) Name() string { return "found" }

func (c FoundCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanFound(pNum) {
		return "", errors.New("can't found at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify which hotel to found")
	}
	corp, err := FindCorp(args[0])
	if err != nil {
		return "", err
	}
	return "", g.FoundCorp(pNum, corp)
}

func (c FoundCommand) Usage(player string, context interface{}) string {
	return `{{b}}found ##{{_b}} to found a corporation on your tile.  Eg. {{b}}found festival{{_b}} or {{b}}found fe{{_b}}`
}

func (g *Game) CanFound(player int) bool {
	return !g.IsFinished() && g.CurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_FOUND_CORP
}
