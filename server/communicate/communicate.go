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
	return email.SendGame(id, g, to, commands, header, initial)
}
