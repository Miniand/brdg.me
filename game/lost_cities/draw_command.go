package lost_cities

import (
	"github.com/beefsack/brdg.me/command"
)

type DrawCommand struct{}

func (d DrawCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("draw", 0, input)
}

func (d DrawCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return g.Players[g.CurrentlyMoving] == player &&
		g.TurnPhase == TURN_PHASE_DRAW && !g.IsFinished()
}

func (d DrawCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return err
	}
	return g.DrawCard(playerNum)
}

func (d DrawCommand) Usage(player string, context interface{}) string {
	return "{{b}}draw{{_b}} to draw a card from the draw pile"
}
