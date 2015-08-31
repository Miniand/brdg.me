package splendor

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type VisitCommand struct{}

func (c VisitCommand) Name() string { return "visit" }

func (c VisitCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("you must specify which noble")
	}
	vNum, err := strconv.Atoi(args[0])
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
