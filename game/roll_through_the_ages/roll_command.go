package roll_through_the_ages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type RollCommand struct{}

func (c RollCommand) Name() string { return "roll" }

func (c RollCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("you must specify at least one dice to roll")
	}
	dice := make([]int, len(args))
	for i, s := range args {
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
	return g.Phase == PhaseRoll && g.RemainingRolls > 0 ||
		g.Phase == PhaseExtraRoll
}

func (g *Game) Roll(player int, diceNum []int) error {
	if !g.CanRoll(player) {
		return errors.New("you can't roll at the moment")
	}
	if len(diceNum) == 0 {
		return errors.New("you must specify which dice to roll")
	}
	if g.Phase == PhaseExtraRoll && len(diceNum) > 1 {
		return errors.New("you may only roll one dice on the extra roll")
	}
	l := len(g.RolledDice)
	for _, n := range diceNum {
		if n < 0 || n > l {
			return fmt.Errorf("dice number must be between 1 and %d", l)
		}
	}
	kept := []int{}
	for i, d := range g.RolledDice {
		if !ContainsInt(i+1, diceNum) {
			kept = append(kept, d)
		}
	}
	rolled := RollN(len(g.RolledDice) - len(kept))
	g.RolledDice = append(rolled, kept...)
	g.LogRoll(rolled, append(kept, g.KeptDice...))
	g.KeepSkulls()
	switch g.Phase {
	case PhaseRoll:
		g.RemainingRolls -= 1
		if g.RemainingRolls == 0 {
			g.NextPhase()
		}
	case PhaseExtraRoll:
		g.NextPhase()
	}
	return nil
}

func (g *Game) NewRoll(n int) {
	g.RolledDice = RollN(n)
	g.LogRoll(g.RolledDice, []int{})
	g.KeptDice = []int{}
	g.KeepSkulls()
}

func (g *Game) KeepSkulls() {
	if len(g.Players) == 1 {
		// You can reroll skulls in single player
		return
	}
	i := 0
	for i < len(g.RolledDice) {
		switch g.RolledDice[i] {
		case DiceSkull:
			g.RolledDice = append(g.RolledDice[:i], g.RolledDice[i+1:]...)
			g.KeptDice = append(g.KeptDice, DiceSkull)
		default:
			i += 1
			continue
		}
	}
	if len(g.RolledDice) == 0 && !(g.Phase == PhaseExtraRoll &&
		g.Boards[g.CurrentPlayer].Developments[DevelopmentLeadership]) {
		g.NextPhase()
	}
}

func (g *Game) LogRoll(newDice, oldDice []int) {
	diceStrings := []string{}
	for _, d := range newDice {
		diceStrings = append(diceStrings, render.Bold(RenderDice(d)))
	}
	for _, d := range oldDice {
		diceStrings = append(diceStrings, RenderDice(d))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s rolled  %s`,
		g.RenderName(g.CurrentPlayer),
		strings.Join(diceStrings, "  "),
	)))
}
