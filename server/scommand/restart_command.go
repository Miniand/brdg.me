package scommand

import (
	"errors"
	"strings"

	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

type RestartCommand struct {
	gameModel *model.GameModel
}

func (rc RestartCommand) Parse(input string) []string {
	return comm.ParseNamedCommand("restart", input)
}

func (rc RestartCommand) CanCall(player string, context interface{}) bool {
	g, ok := context.(game.Playable)
	return ok && g.IsFinished() && rc.gameModel != nil && !rc.gameModel.Restarted
}

func (rc RestartCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if rc.gameModel == nil {
		return "", errors.New("No game was passed in")
	}
	if rc.gameModel.Restarted {
		return "", errors.New("The game has already been restarted")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}
	others := []string{}
	for _, p := range g.PlayerList() {
		if p != player {
			others = append(others, p)
		}
	}
	nc := NewCommand{}
	if _, err := nc.Call(player, nil, []string{
		"new",
		g.Identifier(),
		strings.Join(others, " "),
	}); err != nil {
		return "", err
	}
	rc.gameModel.Restarted = true
	return "The game has been restarted in a new email thread", nil
}

func (rc RestartCommand) Usage(player string, context interface{}) string {
	return "{{b}}restart{{_b}} to restart this game"
}
