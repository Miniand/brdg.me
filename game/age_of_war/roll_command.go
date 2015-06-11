package age_of_war

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (c RollCommand) Parse(input string) []string {
	return command.ParseNamedCommand("roll", input)
}

func (c RollCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanRoll(pNum)
}

func (c RollCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	return "", g.RollForPlayer(pNum)
}

func (c RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll{{_b}} to discard a die and roll the rest"
}

func (g *Game) CanRoll(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) RollForPlayer(player int) error {
	if !g.CanRoll(player) {
		return errors.New("unable to roll right now")
	}
	g.Roll(len(g.CurrentRoll) - 1)
	g.CheckEndOfTurn()
	return nil
}
