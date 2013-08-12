package lost_cities

import (
	"errors"
	"github.com/beefsack/brdg.me/command"
)

type PlayCommand struct{}

func (d PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (d PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.Players[g.CurrentlyMoving] == player &&
		g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
}

func (d PlayCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return errors.New("You must specify a card to play, such as r5")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	c, err := g.ParseCardString(a[0])
	if err != nil {
		return err
	}
	return g.PlayCard(playerNum, c)
}

func (d PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a card, eg. {{b}}play r5{{_b}}"
}
