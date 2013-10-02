package acquire

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseRegexp(fmt.Sprintf("play (%s)", TILE_REGEXP), input)
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
	t, err := ParseTileText(args[1])
	if err != nil {
		return "", err
	}
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.PlayTile(playerNum, t)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a tile to the board.  Eg. {{b}}play 10d{{_b}}"
}
