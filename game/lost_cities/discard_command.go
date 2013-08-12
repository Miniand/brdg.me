package lost_cities

import (
	"errors"
	"github.com/beefsack/brdg.me/command"
)

type DiscardCommand struct{}

func (d DiscardCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("discard", 1, input)
}

func (d DiscardCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.Players[g.CurrentlyMoving] == player &&
		g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
}

func (d DiscardCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return errors.New("You must specify a card to discard, such as r5")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	c, err := g.ParseCardString(a[0])
	if err != nil {
		return err
	}
	return g.DiscardCard(playerNum, c)
}

func (d DiscardCommand) Usage(player string, context interface{}) string {
	return "{{b}}discard ##{{_b}} to discard a card, eg. {{b}}discard r5{{_b}}"
}
