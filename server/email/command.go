package main

import (
	"github.com/beefsack/brdg.me/command"
)

func Commands() []command.Command {
	return []command.Command{
		NewCommand{},
	}
}
