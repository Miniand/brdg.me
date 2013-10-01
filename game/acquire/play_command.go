package acquire

import (
	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("play", 1, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	return !g.IsFinished() && g.CurrentPlayer == playerNum &&
		g.TurnPhase == TURN_PHASE_PLAY_TILE
}

func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	t, err := ParseTileText(args[0])
	if err != nil {
		return "", err
	}
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	err = g.PlayTile(playerNum, t)
	return "", err
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a tile to the board.  Eg. {{b}}play 10d{{_b}}"
}
