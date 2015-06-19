package king_of_tokyo

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type RollCommand struct{}

func (c RollCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("roll", 1, -1, input)
}

func (c RollCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanRoll(pNum)
}

func (c RollCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	dice := make([]int, len(a))
	for i, s := range a {
		d, err := strconv.Atoi(s)
		if err != nil {
			return "", fmt.Errorf("%s is not a number", s)
		}
		dice[i] = d
	}
	return "", g.Roll(pNum, dice)
}

func (c RollCommand) Usage(player string, context interface{}) string {
	return "{{b}}roll # # #{{_b}} to reroll dice, eg. {{b}}roll 1 3 4{{_b}}"
}

func (g *Game) CanRoll(player int) bool {
	if g.CurrentPlayer != player {
		return false
	}
	return g.Phase == PhaseRoll && (g.RemainingRolls > 0 ||
		len(g.ExtraRollable) > 0)
}

func (g *Game) Roll(player int, diceNum []int) error {
	if !g.CanRoll(player) {
		return errors.New("you can't roll at the moment")
	}
	if len(diceNum) == 0 {
		return errors.New("you must specify which dice to roll")
	}
	l := len(g.CurrentRoll)
	for _, n := range diceNum {
		if n < 0 || n > l {
			return fmt.Errorf("dice number must be between 1 and %d", l)
		}
		if g.RemainingRolls <= 0 && !g.ExtraRollable[n-1] {
			return fmt.Errorf("dice number %d is not rerollable", n)
		}
	}
	kept := []int{}
	for i, d := range g.CurrentRoll {
		if !ContainsInt(i+1, diceNum) {
			kept = append(kept, d)
		}
	}
	rolled := RollDice(len(g.CurrentRoll) - len(kept))
	g.LogRoll(player, rolled, kept)
	g.CurrentRoll = append(rolled, kept...)
	switch g.Phase {
	case PhaseRoll:
		g.RemainingRolls -= 1
		if g.RemainingRolls <= 0 {
			g.CheckRollComplete()
		}
	}
	return nil
}
