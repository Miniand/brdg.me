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

func (sc SayCommand) Name() string { return "say" }

func (sc SayCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}

	line, err := input.ReadToEndOfLine()
	if err != nil || line == "" {
		return "", errors.New("please specify something to say")
	}

	message := fmt.Sprintf(
		`{{b}}%s says: %s{{_b}}`,
		render.PlayerNameInPlayers(player, sc.gameModel.PlayerList),
		render.RenderPlain(line),
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
			CommandsForGame(player, sc.gameModel, g),
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

func CanSay(player string, gm *model.GameModel) bool {
	return true
}
