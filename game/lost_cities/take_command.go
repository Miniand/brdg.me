package lost_cities

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type TakeCommand struct{}

func (d TakeCommand) Name() string { return "take" }

func (d TakeCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	if !g.CanTake(playerNum) {
		return "", errors.New("can't take at the moment")
	}
	suitnum := 0
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify a type of card to take, such as r")
	}
	suit := args[0]
	switch suit {
	case "r":
		suitnum = SUIT_RED
	case "y":
		suitnum = SUIT_YELLOW
	case "b":
		suitnum = SUIT_BLUE
	case "w":
		suitnum = SUIT_WHITE
	case "g":
		suitnum = SUIT_GREEN
	default:
		return "", errors.New("could not parse suit")
	}
	return "", g.TakeCard(playerNum, suitnum)
}

func (d TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take #{{_b}} to take a card from a discard pile, eg. {{b}}take r{{_b}}"
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentlyMoving == player &&
		g.TurnPhase == TURN_PHASE_DRAW && !g.IsFinished()
}
