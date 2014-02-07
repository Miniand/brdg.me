package farkle

import (
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type DoneCommand struct{}

func (dc DoneCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("done", 0, input)
}

func (dc DoneCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.Player] && g.TakenThisRoll &&
		!g.IsFinished()
}

func (dc DoneCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	if player != g.Players[g.Player] {
		return "", errors.New("It's not your turn")
	}
	if g.TurnScore == 0 {
		return "", errors.New("You haven't scored anything yet")
	}
	if g.IsFinished() {
		return "", errors.New("The game is already finished")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s took %d points, now on %d",
		render.PlayerName(g.Player, g.Players[g.Player]),
		g.TurnScore, g.Scores[g.Player]+g.TurnScore)))
	g.Scores[g.Player] = g.Scores[g.Player] + g.TurnScore
	g.NextPlayer()
	return "", nil
}

func (dc DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to take the points and finish your turn"
}
