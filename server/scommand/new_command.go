package scommand

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/model"
)

const (
	MsgTypeInvite = "invite"
)

type NewCommand struct{}

func (nc NewCommand) Name() string { return "name" }

func (nc NewCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	args, err := input.ReadLineArgs()
	if err != nil || len(args) < 2 {
		errors.New("Could not find game name and players")
	}
	g := game.RawCollection()[args[0]]
	players := append([]string{player}, args[1:]...)
	gm, err := model.StartNewGame(g, players)
	if err != nil {
		return "", err
	}
	return "", communicate.Game(
		g,
		gm,
		gm.PlayerList,
		CommandsForGame(player, gm, g),
		"You have been invited by "+player+" to play "+g.Name()+"!",
		MsgTypeInvite,
		true,
	)
}

func (nc NewCommand) Usage(player string, context interface{}) string {
	return "{{b}}new (game ID) (players...){{_b}} start a new game with friends"
}

func CanNew(player string, gm *model.GameModel) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || ok && u.Unsubscribed {
		return false
	}
	return gm == nil || gm.IsFinished
}
