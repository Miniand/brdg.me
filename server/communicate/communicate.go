package communicate

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/server/email"
)

func Game(
	id string,
	g game.Playable,
	to []string,
	commands []command.Command,
	header string,
	initial bool,
) error {
	failed := wsSendGameMulti(to, id, header, g)
	if len(failed) == 0 {
		return nil
	}
	emailTo := []string{}
	for p, _ := range failed {
		emailTo = append(emailTo, p)
	}
	return email.SendGame(id, g, emailTo, commands, header, initial)
}

func GameUpdate(id string, g game.Playable, to []string, text string) {
	wsSendGameMulti(to, id, text, g)
}
