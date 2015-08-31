package zombie_dice

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (c RollCommand) Name() string { return "roll" }

func (c RollCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	if !g.CanRoll(pNum) {
		return "", errors.New("cannot roll at the moment")
	}
	g.Roll()
	return "", nil
}

func (c RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll{{_b}} to push your luck and roll the dice"
}

func (g *Game) CanRoll(player int) bool {
	return g.CurrentTurn == player
}
