package scommand

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
)

type ConcedeYesCommand struct {
	gameModel *model.GameModel
}

func (c ConcedeYesCommand) Parse(input string) []string {
	return command.ParseNamedCommand("yes", input)
}

func (c ConcedeYesCommand) CanCall(player string, context interface{}) bool {
	return CanConcedeVote(player, c.gameModel)
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
	g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
		"%s voted {{b}}yes{{_b}}",
		render.PlayerNameInPlayers(player, c.gameModel.PlayerList),
	)))
	if len(c.gameModel.RemainingConcedeVotePlayers()) == 0 {
		PassConcedeVote(c.gameModel, g)
	}
	return "", nil
}

func (c ConcedeYesCommand) Usage(player string, context interface{}) string {
	return "{{b}}yes{{_b}} to vote to concede"
}
