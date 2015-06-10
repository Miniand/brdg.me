package age_of_war

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type LineCommand struct{}

func (c LineCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("line", 1, input)
}

func (c LineCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanLine(pNum)
}

func (c LineCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) != 1 {
		return "", errors.New("you must specify a line to complete")
	}

	line := 0
	return "", g.Line(pNum, line)
}

func (c LineCommand) Usage(player string, context interface{}) string {
	return "{{b}}line #{{_b}} to complete a line in the castle you are attacking, eg. {{b}}line 2{{_b}}"
}

func (g *Game) CanLine(player int) bool {
	return g.CurrentPlayer == player && g.CurrentlyAttacking != -1
}

func (g *Game) Line(player, line int) error {
	if !g.CanLine(player) {
		return errors.New("unable to complete a line right now")
	}
	if line < 0 || line >= len(Castles[g.CurrentlyAttacking].Lines) {
		return errors.New("that is not a valid line")
	}
	return errors.New("not implemented")
}
