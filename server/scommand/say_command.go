package scommand

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/model"
)

const (
	MsgTypeSay = "say"
)

type SayCommand struct {
	gameModel *model.GameModel
}

func (sc SayCommand) Parse(input string) []string {
	return command.ParseRegexp(`say ([^\r\n]+)$`, input)
}

func (sc SayCommand) CanCall(player string, context interface{}) bool {
	_, ok := context.(game.Playable)
	return ok
}

func (sc SayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}
	message := fmt.Sprintf(
		`{{b}}%s says: %s{{_b}}`,
		render.PlayerNameInPlayers(player, sc.gameModel.PlayerList),
		render.RenderPlain(args[1]),
	)
	g.GameLog().Add(log.NewPublicMessage(message))
	if sc.gameModel.IsFinished {
		// Just send it out to everyone.
		otherPlayers := []string{}
		for _, p := range sc.gameModel.PlayerList {
			if p != player {
				otherPlayers = append(otherPlayers, p)
			}
		}
		communicate.Game(
			g,
			sc.gameModel,
			otherPlayers,
			CommandsForGame(sc.gameModel, g),
			message,
			MsgTypeSay,
			false,
		)
	}
	return "", nil
}

func (sc SayCommand) Usage(player string, context interface{}) string {
	return "{{b}}say ##{{_b}} to send a message to the other players, eg. {{b}}say hello!{{_b}}"
}
