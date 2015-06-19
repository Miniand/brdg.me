package king_of_tokyo

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
	if !ok {
		return false
	}
	return g.CanKeep(pNum)
}

func (c KeepCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	return "", g.Keep(pNum)
}

func (c KeepCommand) Usage(player string, context interface{}) string {
	return "{{b}}keep{{_b}} to keep the current roll"
}

func (g *Game) CanKeep(player int) bool {
	if g.CurrentPlayer != player {
		return false
	}
	return g.Phase == PhaseRoll
}

func (g *Game) Keep(player int) error {
	if !g.CanKeep(player) {
		return errors.New("you can't keep at the moment")
	}
	g.NextPhase()
	return nil
}
