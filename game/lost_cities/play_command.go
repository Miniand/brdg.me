package lost_cities

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (d PlayCommand) Name() string { return "play" }

func (d PlayCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	if !g.CanPlay(playerNum) {
		return "", errors.New("can't play at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a card to play, such as r5")
	}
	c, err := g.ParseCardString(args[0])
	if err != nil {
		return "", err
	}
	return "", g.PlayCard(playerNum, c)
}

func (d PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a card, eg. {{b}}play r5{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentlyMoving == player &&
		g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
}
