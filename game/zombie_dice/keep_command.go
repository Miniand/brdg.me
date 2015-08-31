package zombie_dice

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type KeepCommand struct{}

func (c KeepCommand) Name() string { return "keep" }

func (c KeepCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	if !g.CanKeep(pNum) {
		return "", errors.New("cannot keep at the moment")
	}
	g.Keep()
	return "", nil
}

func (c KeepCommand) Usage(player string, context interface{}) string {
	return "{{b}}keep{{_b}} to be a coward and keep your brains"
}

func (g *Game) CanKeep(player int) bool {
	return g.CurrentTurn == player
}
