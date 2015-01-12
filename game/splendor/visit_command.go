package splendor

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type VisitCommand struct{}

func (c VisitCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("visit", 1, input)
}

func (c VisitCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	return err != nil && g.CanVisit(pNum)
}

func (c VisitCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which noble")
	}
	vNum, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Visit(pNum, vNum-1)
}

func (c VisitCommand) Usage(player string, context interface{}) string {
	return "{{b}}visit #{{_b}} to visit a noble, eg. {{b}}visit 2{{_b}}"
}

func (g *Game) CanVisit(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseVisit
}

func (g *Game) Visit(player, noble int) error {
	if !g.CanVisit(player) {
		return errors.New("unable to visit right now")
	}
	if noble < 0 || noble > len(g.Nobles) {
		return errors.New("that is not a valid noble number")
	}
	g.PlayerBoards[player].Nobles = append(g.PlayerBoards[player].Nobles,
		g.Nobles[noble])
	g.Nobles = append(g.Nobles[:noble], g.Nobles[noble+1:]...)
	g.NextPhase()
	return nil
}
