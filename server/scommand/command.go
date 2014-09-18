package scommand

import (
	comm "github.com/Miniand/brdg.me/command"
)

func Commands(gameId string) []comm.Command {
	return []comm.Command{
		PokeCommand{
			gameId: gameId,
		},
		SayCommand{
			gameId: gameId,
		},
		NewCommand{},
		RestartCommand{},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
