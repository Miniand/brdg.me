package communicate

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/model"
)

func Game(
	g game.Playable,
	gm *model.GameModel,
	to []string,
	commands []command.Command,
	header string,
	headerType string,
	initial bool,
) error {
	failed := wsSendGameMulti(to, header, headerType, g, gm)
	if len(failed) == 0 {
		return nil
	}
	emailTo := []string{}
	for p, _ := range failed {
		emailTo = append(emailTo, p)
	}
	return email.SendGame(g, gm, emailTo, commands, header, initial)
}

func GameUpdate(g game.Playable, gm *model.GameModel, to []string, text, msgType string) {
	wsSendGameMulti(to, text, msgType, g, gm)
}
