package main

import (
	"github.com/Miniand/brdg.me/command"
)

func Commands(gameId interface{}) []command.Command {
	return []command.Command{
		PokeCommand{
			gameId: gameId,
		},
		NewCommand{},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
