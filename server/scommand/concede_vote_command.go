package scommand

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
)

type ConcedeVoteCommand struct {
	gameModel *model.GameModel
}

func (c ConcedeVoteCommand) Name() string { return "vote" }

func (c ConcedeVoteCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	if !CanConcedeVote(player, c.gameModel) {
		return "", errors.New("you can't vote at the moment")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("no game was passed in")
	}

	args, err := input.ReadLineArgs()
	if err != nil || len(args) != 1 {
		return "", errors.New("please specify 'yes' or 'no'")
	}
	switch strings.ToLower(args[0]) {
	case "yes":
		g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
			"%s voted {{b}}yes{{_b}}",
			render.PlayerNameInPlayers(player, c.gameModel.PlayerList),
		)))
		PassConcedeVote(c.gameModel, g)
	case "no":
		g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
			"%s voted {{b}}no{{_b}}",
			render.PlayerNameInPlayers(player, c.gameModel.PlayerList),
		)))
		FailConcedeVote(c.gameModel, g)
	default:
		return "", errors.New("please specify 'yes' or 'no'")
	}
	return "", nil
}

func (c ConcedeVoteCommand) Usage(player string, context interface{}) string {
	return "{{b}}vote yes/no{{_b}} to vote for or against the concede"
}
