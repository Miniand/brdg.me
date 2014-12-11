package scommand

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/model"
)

type ConcedeYesCommand struct {
	gameModel *model.GameModel
}

func (c ConcedeYesCommand) Parse(input string) []string {
	return command.ParseNamedCommand("yes", input)
}

func (c ConcedeYesCommand) CanCall(player string, context interface{}) bool {
	return c.gameModel != nil && CanConcedeVote(player, c.gameModel)
}

func (c ConcedeYesCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if !c.CanCall(player, context) {
		return "", errors.New("you can't vote at the moment")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("no game was passed in")
	}

	c.gameModel.ConcedeVote[player] = true
	if len(RemainingConcedeVotePlayers(c.gameModel)) == 0 {
		PassConcedeVote(c.gameModel, g)
	}
	return "", nil
}

func (c ConcedeYesCommand) Usage(player string, context interface{}) string {
	return "{{b}}yes{{_b}} to vote to concede"
}
