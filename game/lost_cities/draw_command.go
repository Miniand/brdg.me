package lost_cities

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type DrawCommand struct{}

func (d DrawCommand) Name() string { return "draw" }

func (d DrawCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	if !g.CanDraw(playerNum) {
		return "", errors.New("can't draw at the moment")
	}
	return "", g.DrawCard(playerNum)
}

func (d DrawCommand) Usage(player string, context interface{}) string {
	return "{{b}}draw{{_b}} to draw a card from the draw pile"
}

func (g *Game) CanDraw(player int) bool {
	return g.CurrentlyMoving == player &&
		g.TurnPhase == TURN_PHASE_DRAW && !g.IsFinished()
}
