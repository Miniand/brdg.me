package splendor

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type VisitCommand struct{}

func (c VisitCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("visit", 1, input)
}

func (c VisitCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	return found && g.CanVisit(pNum)
}

func (c VisitCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
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
	if noble < 0 || noble >= len(g.Nobles) {
		return errors.New("that is not a valid noble number")
	}
	g.PlayerBoards[player].Nobles = append(g.PlayerBoards[player].Nobles,
		g.Nobles[noble])
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s was visited by %s",
		g.RenderName(player),
		RenderNoble(g.Nobles[noble]),
	)))
	g.Nobles = append(g.Nobles[:noble], g.Nobles[noble+1:]...)
	g.NextPhase()
	return nil
}
