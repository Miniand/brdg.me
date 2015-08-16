package red7

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DoneCommand struct{}

func (c DoneCommand) Name() string { return "done" }

func (c DoneCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	return "", g.Done(pNum)
}

func (c DoneCommand) Usage(player string, context interface{}) string {
	return "{{b}}done{{_b}} to finish your turn, you will be eliminated if you haven't played or discarded a card or if you aren't the leader"
}

func (g *Game) CanDone(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Done(player int) error {
	if !g.CanDone(player) {
		return errors.New("you can't done at the moment")
	}
	if !g.HasPlayed {
		g.Eliminate(g.CurrentPlayer, " for not playing or discarding")
	}
	g.EndTurn()
	return nil
}
