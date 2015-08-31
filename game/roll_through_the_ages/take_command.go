package roll_through_the_ages

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

const (
	TakeFood = iota
	TakeWorkers
)

var TakeMap = map[int]string{
	TakeFood:    "food",
	TakeWorkers: "workers",
}

type TakeCommand struct{}

func (c TakeCommand) Name() string { return "take" }

func (c TakeCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}

	actions := []int{}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("you must specify at least one thing to take")
	}
	for _, a := range args {
		action, err := helper.MatchStringInStringMap(a, TakeMap)
		if err != nil {
			return "", err
		}
		actions = append(actions, action)
	}

	return "", g.Take(pNum, actions)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return fmt.Sprintf(
		"{{b}}take # # #{{_b}} to take food or workers, one for each %s dice, eg. for two dice, {{b}}take food workers{{_b}}",
		RenderDice(DiceFoodOrWorkers),
	)
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseCollect
}

func (g *Game) Take(player int, actions []int) error {
	if !g.CanTake(player) {
		return errors.New("you can't take at the moment")
	}
	numDice := 0
	for _, d := range g.KeptDice {
		if d == DiceFoodOrWorkers {
			numDice += 1
		}
	}
	if l := len(actions); l != numDice {
		return fmt.Errorf(
			"you must specify %d take actions after the take command",
			l,
		)
	}

	cp := g.CurrentPlayer
	for _, a := range actions {
		switch a {
		case TakeFood:
			g.Boards[cp].Food += 2 + g.Boards[cp].FoodModifier()
		case TakeWorkers:
			g.RemainingWorkers += 2 + g.Boards[cp].WorkerModifier()
		default:
			return errors.New("could not understand action")
		}
	}

	g.NextPhase()
	return nil
}
