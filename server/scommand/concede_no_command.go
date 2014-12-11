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

type ConcedeNoCommand struct {
	gameModel *model.GameModel
}

func (c ConcedeNoCommand) Parse(input string) []string {
	return command.ParseNamedCommand("no", input)
}

func (c ConcedeNoCommand) CanCall(player string, context interface{}) bool {
	return CanConcedeVote(player, c.gameModel)
}

func (c ConcedeNoCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if !c.CanCall(player, context) {
		return "", errors.New("you can't vote at the moment")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("no game was passed in")
	}

	g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
		"%s voted {{b}}no{{_b}}",
		render.PlayerNameInPlayers(player, c.gameModel.PlayerList),
	)))
	FailConcedeVote(c.gameModel, g)
	return "", nil
}

func (c ConcedeNoCommand) Usage(player string, context interface{}) string {
	return "{{b}}no{{_b}} to vote against conceeding"
}
