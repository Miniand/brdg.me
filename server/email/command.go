package email

import (
	"github.com/Miniand/brdg.me/command"
)

func Commands(gameId string) []command.Command {
	return []command.Command{
		PokeCommand{
			gameId: gameId,
		},
		NewCommand{},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
