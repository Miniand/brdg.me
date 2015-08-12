package farkle

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (rc RollCommand) Name() string { return "roll" }

func (rc RollCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("can't find player")
	}
	if !g.CanRoll(pNum) {
		return "", errors.New("can't play at the moment")
	}
	if len(g.RemainingDice) > 0 {
		g.Roll(len(g.RemainingDice))
	} else {
		g.Roll(6)
	}
	g.TakenThisRoll = false
	return "", nil
}

func (rc RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll{{_b}} to roll remaining dice"
}

func (g *Game) CanRoll(player int) bool {
	return player == g.Player && g.TakenThisRoll && !g.IsFinished()
}
