package zombie_dice

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type KeepCommand struct{}

func (c KeepCommand) Parse(input string) []string {
	return command.ParseNamedCommand("keep", input)
}

func (c KeepCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanKeep(pNum)
}

func (c KeepCommand) Call(player string, context interface{},
	args []string) (string, error) {
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
