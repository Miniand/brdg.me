package roll_through_the_ages

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

type KeepCommand struct{}

func (c KeepCommand) Parse(input string) []string {
	return command.ParseNamedCommand("keep", input)
}

func (c KeepCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return g.CanKeep(pNum)
}

func (c KeepCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Keep(pNum)
}

func (c KeepCommand) Usage(player string, context interface{}) string {
	return "{{b}}keep{{_b}} to keep the current dice"
}

func (g *Game) CanKeep(player int) bool {
	return g.CanRoll(player)
}

func (g *Game) Keep(player int) error {
	if !g.CanKeep(player) {
		return errors.New("you can't keep at the moment")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s kept the current dice`,
		g.RenderName(player),
	)))
	g.CollectPhase()
	return nil
}
