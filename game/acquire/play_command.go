package acquire

import (
	"errors"
	"log"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Name() string { return "play" }

func (c PlayCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	if !g.CanPlay(pNum) {
		return "", errors.New("can't play at the moment")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		log.Println(args, err)
		return "", errors.New("please specify which tile to play")
	}
	t, err := ParseTileText(args[0])
	if err != nil {
		return "", err
	}
	return "", g.PlayTile(pNum, t)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play ##{{_b}} to play a tile to the board.  Eg. {{b}}play 10d{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return !g.IsFinished() && g.CurrentPlayer == player &&
		g.TurnPhase == TURN_PHASE_PLAY_TILE
}
