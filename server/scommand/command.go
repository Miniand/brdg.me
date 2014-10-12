package scommand

import (
	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/model"
)

func Commands(gm *model.GameModel) []comm.Command {
	return []comm.Command{
		PokeCommand{
			gameModel: gm,
		},
		SayCommand{
			gameModel: gm,
		},
		NewCommand{},
		RestartCommand{
			gameModel: gm,
		},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
