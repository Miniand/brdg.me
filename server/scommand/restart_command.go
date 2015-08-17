package scommand

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

type RestartCommand struct {
	gameModel *model.GameModel
}

func (rc RestartCommand) Name() string { return "restart" }

func (rc RestartCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	if rc.gameModel == nil {
		return "", errors.New("no game was passed in")
	}
	if rc.gameModel.Restarted {
		return "", errors.New("the game has already been restarted")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}
	others := []string{}
	for _, p := range rc.gameModel.PlayerList {
		if p != player {
			others = append(others, p)
		}
	}
	nc := NewCommand{}
	if _, err := nc.Call(player, nil, command.NewParserString(fmt.Sprintf(
		"%s %s",
		g.Identifier(),
		strings.Join(others, " "),
	))); err != nil {
		return "", err
	}
	rc.gameModel.Restarted = true
	return "The game has been restarted", nil
}

func (rc RestartCommand) Usage(player string, context interface{}) string {
	return "{{b}}restart{{_b}} to restart this game"
}

func CanRestart(player string, gm *model.GameModel) bool {
	return gm.IsFinished && !gm.Restarted
}
