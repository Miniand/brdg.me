package main

import (
	"github.com/Miniand/brdg.me/command"
)

func Commands() []command.Command {
	return []command.Command{
		PokeCommand{},
		NewCommand{},
		UnsubscribeCommand{},
		SubscribeCommand{},
	}
}
