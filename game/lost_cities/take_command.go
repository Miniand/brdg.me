package lost_cities

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
)

type TakeCommand struct{}

func (d TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 1, input)
}

func (d TakeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.Players[g.CurrentlyMoving] == player &&
		g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
}

func (d TakeCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return errors.New("You must specify a type of card to take, such as r")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	return g.TakeCard(playerNum, SUIT_RED)
}

func (d TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take #{{_b}} to take a card from a discard pile, eg. {{b}}take r{{_b}}"
}
