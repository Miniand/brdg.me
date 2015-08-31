package lost_cities

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DiscardCommand struct{}

func (d DiscardCommand) Name() string { return "discard" }

func (d DiscardCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a card to discard, such as r5")
	}
	if !g.CanDiscard(playerNum) {
		return "", errors.New("can't discard at the moment")
	}
	c, err := g.ParseCardString(args[0])
	if err != nil {
		return "", err
	}
	return "", g.DiscardCard(playerNum, c)
}

func (d DiscardCommand) Usage(player string, context interface{}) string {
	return "{{b}}discard ##{{_b}} to discard a card, eg. {{b}}discard r5{{_b}}"
}

func (g *Game) CanDiscard(player int) bool {
	return g.CurrentlyMoving == player &&
		g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
}
