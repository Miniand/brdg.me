package farkle

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type DoneCommand struct{}

func (dc DoneCommand) Name() string { return "done" }

func (dc DoneCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("cannot find player")
	}
	if !g.CanDone(pNum) {
		return "", errors.New("can't call done at the moment")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took {{b}}%d{{_b}} points, now on {{b}}%d{{_b}}",
		render.PlayerName(g.Player, g.Players[g.Player]),
		g.TurnScore,
		g.Scores[pNum]+g.TurnScore,
	)))
	g.Scores[pNum] = g.Scores[pNum] + g.TurnScore
	g.NextPlayer()
	return "", nil
}

func (dc DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to take the points and finish your turn"
}

func (g *Game) CanDone(player int) bool {
	return player == g.Player && g.TakenThisRoll && !g.IsFinished()
}
